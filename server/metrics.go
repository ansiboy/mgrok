package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mgrok/log"
	"net"
	"net/http"
	"os"
	"time"

	gometrics "github.com/rcrowley/go-metrics"
)

var metrics Metrics

func init() {
	keenAPIKey := os.Getenv("KEEN_API_KEY")

	if keenAPIKey != "" {
		metrics = newKeenIoMetrics(60 * time.Second)
	} else {
		metrics = newLocalMetrics(30 * time.Second)
	}
}

// Metrics metrics
type Metrics interface {
	log.Logger
	openConnection(*Tunnel, net.Conn)
	closeConnection(*Tunnel, net.Conn, time.Time, int64, int64)
	openTunnel(*Tunnel)
	closeTunnel(*Tunnel)
}

// LocalMetrics local metrics
type LocalMetrics struct {
	log.Logger
	reportInterval time.Duration
	windowsCounter gometrics.Counter
	linuxCounter   gometrics.Counter
	osxCounter     gometrics.Counter
	otherCounter   gometrics.Counter

	tunnelMeter        gometrics.Meter
	tcpTunnelMeter     gometrics.Meter
	httpTunnelMeter    gometrics.Meter
	connMeter          gometrics.Meter
	lostHeartbeatMeter gometrics.Meter

	connTimer gometrics.Timer

	bytesInCount  gometrics.Counter
	bytesOutCount gometrics.Counter

	/*
	   tunnelGauge gometrics.Gauge
	   tcpTunnelGauge gometrics.Gauge
	   connGauge gometrics.Gauge
	*/
}

func newLocalMetrics(reportInterval time.Duration) *LocalMetrics {
	metrics := LocalMetrics{
		Logger:         log.NewPrefixLogger("metrics"),
		reportInterval: reportInterval,
		windowsCounter: gometrics.NewCounter(),
		linuxCounter:   gometrics.NewCounter(),
		osxCounter:     gometrics.NewCounter(),
		otherCounter:   gometrics.NewCounter(),

		tunnelMeter:        gometrics.NewMeter(),
		tcpTunnelMeter:     gometrics.NewMeter(),
		httpTunnelMeter:    gometrics.NewMeter(),
		connMeter:          gometrics.NewMeter(),
		lostHeartbeatMeter: gometrics.NewMeter(),

		connTimer: gometrics.NewTimer(),

		bytesInCount:  gometrics.NewCounter(),
		bytesOutCount: gometrics.NewCounter(),

		/*
		   metrics.tunnelGauge = gometrics.NewGauge(),
		   metrics.tcpTunnelGauge = gometrics.NewGauge(),
		   metrics.connGauge = gometrics.NewGauge(),
		*/
	}

	go metrics.report()

	return &metrics
}

func (m *LocalMetrics) openTunnel(t *Tunnel) {
	m.tunnelMeter.Mark(1)

	switch t.ctl.auth.OS {
	case "windows":
		m.windowsCounter.Inc(1)
	case "linux":
		m.linuxCounter.Inc(1)
	case "darwin":
		m.osxCounter.Inc(1)
	default:
		m.otherCounter.Inc(1)
	}

	switch t.req.Protocol {
	case "tcp":
		m.tcpTunnelMeter.Mark(1)
	case "http":
		m.httpTunnelMeter.Mark(1)
	}
}

func (m *LocalMetrics) closeTunnel(t *Tunnel) {
}

func (m *LocalMetrics) openConnection(t *Tunnel, c net.Conn) {
	m.connMeter.Mark(1)
}

func (m *LocalMetrics) closeConnection(t *Tunnel, c net.Conn, start time.Time, bytesIn, bytesOut int64) {
	m.bytesInCount.Inc(bytesIn)
	m.bytesOutCount.Inc(bytesOut)
}

func (m *LocalMetrics) report() {
	m.Info("Reporting every %d seconds", int(m.reportInterval.Seconds()))

	for {
		time.Sleep(m.reportInterval)
		buffer, err := json.Marshal(map[string]interface{}{
			"windows":               m.windowsCounter.Count(),
			"linux":                 m.linuxCounter.Count(),
			"osx":                   m.osxCounter.Count(),
			"other":                 m.otherCounter.Count(),
			"httpTunnelMeter.count": m.httpTunnelMeter.Count(),
			"tcpTunnelMeter.count":  m.tcpTunnelMeter.Count(),
			"tunnelMeter.count":     m.tunnelMeter.Count(),
			"tunnelMeter.m1":        m.tunnelMeter.Rate1(),
			"connMeter.count":       m.connMeter.Count(),
			"connMeter.m1":          m.connMeter.Rate1(),
			"bytesIn.count":         m.bytesInCount.Count(),
			"bytesOut.count":        m.bytesOutCount.Count(),
		})

		if err != nil {
			m.Error("Failed to serialize metrics: %v", err)
			continue
		}

		m.Info("Reporting: %s", buffer)
	}
}

// KeenIoMetric KeenIoMetric
type KeenIoMetric struct {
	Collection string
	Event      interface{}
}

// KeenIoMetrics KeenIoMetrics
type KeenIoMetrics struct {
	log.Logger
	APIKey       string
	ProjectToken string
	HTTPClient   http.Client
	Metrics      chan *KeenIoMetric
}

func newKeenIoMetrics(batchInterval time.Duration) *KeenIoMetrics {
	k := &KeenIoMetrics{
		Logger:       log.NewPrefixLogger("metrics"),
		APIKey:       os.Getenv("KEEN_API_KEY"),
		ProjectToken: os.Getenv("KEEN_PROJECT_TOKEN"),
		Metrics:      make(chan *KeenIoMetric, 1000),
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				k.Error("KeenIoMetrics failed: %v", r)
			}
		}()

		batch := make(map[string][]interface{})
		batchTimer := time.Tick(batchInterval)

		for {
			select {
			case m := <-k.Metrics:
				list, ok := batch[m.Collection]
				if !ok {
					list = make([]interface{}, 0)
				}
				batch[m.Collection] = append(list, m.Event)

			case <-batchTimer:
				// no metrics to report
				if len(batch) == 0 {
					continue
				}

				payload, err := json.Marshal(batch)
				if err != nil {
					k.Error("Failed to serialize metrics payload: %v, %v", batch, err)
				} else {
					for key, val := range batch {
						k.Debug("Reporting %d metrics for %s", len(val), key)
					}

					k.authedRequest("POST", "/events", bytes.NewReader(payload))
				}
				batch = make(map[string][]interface{})
			}
		}
	}()

	return k
}

func (k *KeenIoMetrics) authedRequest(method, path string, body *bytes.Reader) (resp *http.Response, err error) {
	path = fmt.Sprintf("https://api.keen.io/3.0/projects/%s%s", k.ProjectToken, path)
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", k.APIKey)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(body.Len())
	}

	requestStartAt := time.Now()
	resp, err = k.HTTPClient.Do(req)

	if err != nil {
		k.Error("Failed to send metric event to keen.io %v", err)
	} else {
		k.Info("keen.io processed request in %f sec", time.Since(requestStartAt).Seconds())
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			bytes, _ := ioutil.ReadAll(resp.Body)
			k.Error("Got %v response from keen.io: %s", resp.StatusCode, bytes)
		}
	}

	return
}

func (k *KeenIoMetrics) openConnection(t *Tunnel, c net.Conn) {
}

func (k *KeenIoMetrics) closeConnection(t *Tunnel, c net.Conn, start time.Time, in, out int64) {
	event := struct {
		Keen               keenStruct `json:"keen"`
		OS                 string
		ClientID           string
		Protocol           string
		URL                string
		User               string
		Version            string
		Reason             string
		HTTPAuth           bool
		Subdomain          bool
		TunnelDuration     float64
		ConnectionDuration float64
		BytesIn            int64
		BytesOut           int64
	}{
		Keen: keenStruct{
			Timestamp: start.UTC().Format("2006-01-02T15:04:05.000Z"),
		},
		OS:                 t.ctl.auth.OS,
		ClientID:           t.ctl.id,
		Protocol:           t.req.Protocol,
		URL:                t.url,
		User:               t.ctl.auth.User,
		Version:            t.ctl.auth.MmVersion,
		HTTPAuth:           t.req.HTTPAuth != "",
		Subdomain:          t.req.Subdomain != "",
		TunnelDuration:     time.Since(t.start).Seconds(),
		ConnectionDuration: time.Since(start).Seconds(),
		BytesIn:            in,
		BytesOut:           out,
	}

	k.Metrics <- &KeenIoMetric{Collection: "CloseConnection", Event: event}
}

func (k *KeenIoMetrics) openTunnel(t *Tunnel) {
}

type keenStruct struct {
	Timestamp string `json:"timestamp"`
}

func (k *KeenIoMetrics) closeTunnel(t *Tunnel) {
	event := struct {
		Keen      keenStruct `json:"keen"`
		OS        string
		ClientID  string
		Protocol  string
		URL       string
		User      string
		Version   string
		Reason    string
		Duration  float64
		HTTPAuth  bool
		Subdomain bool
	}{
		Keen: keenStruct{
			Timestamp: t.start.UTC().Format("2006-01-02T15:04:05.000Z"),
		},
		OS:       t.ctl.auth.OS,
		ClientID: t.ctl.id,
		Protocol: t.req.Protocol,
		URL:      t.url,
		User:     t.ctl.auth.User,
		Version:  t.ctl.auth.MmVersion,
		//Reason: reason,
		Duration:  time.Since(t.start).Seconds(),
		HTTPAuth:  t.req.HTTPAuth != "",
		Subdomain: t.req.Subdomain != "",
	}

	k.Metrics <- &KeenIoMetric{Collection: "CloseTunnel", Event: event}
}
