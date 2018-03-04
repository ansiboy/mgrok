package client

import (
	"fmt"
	"math"
	"mgrok/conn"
	"mgrok/log"
	"mgrok/msg"
	"mgrok/util"
	"mgrok/version"
	"net"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

const (
	defaultServerAddr   = "t.mgrok.cn:4443"
	defaultInspectAddr  = "127.0.0.1:4040"
	pingInterval        = 20 * time.Second
	maxPongLatency      = 15 * time.Second
	updateCheckInterval = 6 * time.Hour
	badGateway          = `<html>
<body style="background-color: #97a8b9">
    <div style="margin:auto; width:400px;padding: 20px 60px; background-color: #D3D3D3; border: 5px solid maroon;">
        <h2>Tunnel %s unavailable</h2>
        <p>Unable to initiate connection to <strong>%s</strong>. A web server must be running on port <strong>%s</strong> to complete the tunnel.</p>
`
)

// Model client model
type Model struct {
	log.Logger
	id             string
	tunnels        map[string]Tunnel
	serverVersion  string
	metrics        *Metrics
	updateStatus   UpdateStatus
	connStatus     ConnStatus
	serverAddr     string
	proxyURL       string
	authToken      string
	tunnelConfig   map[string]*TunnelConfiguration
	configPath     string
	updateCallback func(c *Model)
}

func newClientModel(config *Configuration) *Model {

	m := &Model{
		Logger: log.NewPrefixLogger("client"),

		// server address
		serverAddr: config.ServerAddr,

		// proxy address
		proxyURL: config.HTTPProxy,

		// auth token
		authToken: config.AuthToken,

		// connection status
		connStatus: ConnConnecting,

		// update status
		updateStatus: UpdateNone,

		// metrics
		metrics: NewClientMetrics(),

		// open tunnels
		tunnels: make(map[string]Tunnel),

		// tunnel configuration
		tunnelConfig: config.Tunnels,

		// config path
		configPath: config.Path,
	}

	return m
}

// server name in release builds is the host part of the server address
func serverName(addr string) string {
	host, _, err := net.SplitHostPort(addr)

	// should never panic because the config parser calls SplitHostPort first
	if err != nil {
		panic(err)
	}

	return host
}

// GetClientVersion get client version
func (c Model) GetClientVersion() string { return version.MajorMinor() }

// GetServerVersion get server version
func (c Model) GetServerVersion() string { return c.serverVersion }

// GetTunnels get tunnels
func (c Model) GetTunnels() []Tunnel {
	tunnels := make([]Tunnel, 0)
	for _, t := range c.tunnels {
		tunnels = append(tunnels, t)
	}
	return tunnels
}

// GetConnStatus get connection status
func (c Model) GetConnStatus() ConnStatus { return c.connStatus }

// GetUpdateStatus client update status
func (c Model) GetUpdateStatus() UpdateStatus { return c.updateStatus }

// GetConnectionMetrics connection metrics
func (c Model) GetConnectionMetrics() (metrics.Meter, metrics.Timer) {
	return c.metrics.connMeter, c.metrics.connTimer
}

// GetBytesInMetrics bytes in metrics
func (c Model) GetBytesInMetrics() (metrics.Counter, metrics.Histogram) {
	return c.metrics.bytesInCount, c.metrics.bytesIn
}

// GetBytesOutMetrics bytes out metrics
func (c Model) GetBytesOutMetrics() (metrics.Counter, metrics.Histogram) {
	return c.metrics.bytesOutCount, c.metrics.bytesOut
}

// SetUpdateStatus set update status
func (c Model) SetUpdateStatus(updateStatus UpdateStatus) {
	c.updateStatus = updateStatus
	c.update()
}

func (c *Model) update() {
	if c.updateCallback != nil {
		c.updateCallback(c)
	}

}

// Run run
func (c *Model) Run() {
	// how long we should wait before we reconnect
	maxWait := 30 * time.Second
	wait := 1 * time.Second

	for {
		// run the control channel
		c.control()

		// control only returns when a failure has occurred, so we're going to try to reconnect
		if c.connStatus == ConnOnline {
			wait = 1 * time.Second
		}

		log.Info("Waiting %d seconds before reconnecting", int(wait.Seconds()))
		time.Sleep(wait)
		// exponentially increase wait time
		wait = 2 * wait
		wait = time.Duration(math.Min(float64(wait), float64(maxWait)))
		c.connStatus = ConnReconnecting
		c.update()
	}
}

// Establishes and manages a tunnel control connection with the server
func (c *Model) control() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("control recovering from failure %v", r)
		}
	}()

	// establish control channel
	var (
		ctlConn net.Conn
		err     error
	)
	if c.proxyURL == "" {
		// simple non-proxied case, just connect to the server
		ctlConn, err = conn.Dial(c.serverAddr, "ctl")
	} else {
		ctlConn, err = conn.DialHttpProxy(c.proxyURL, c.serverAddr, "ctl")
	}
	if err != nil {
		panic(err)
	}
	defer func() {
		ctlConn.Close()
	}()

	// authenticate with the server
	auth := &msg.Auth{
		ClientId:  c.id,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   version.Proto,
		MmVersion: version.MajorMinor(),
		User:      c.authToken,
	}

	if err = msg.WriteMsg(ctlConn, auth); err != nil {
		panic(err)
	}

	// wait for the server to authenticate us
	var authResp msg.AuthResp
	if err = msg.ReadMsgInto(ctlConn, &authResp); err != nil {
		panic(err)
	}

	if authResp.Error != "" {
		emsg := fmt.Sprintf("Failed to authenticate to server: %s", authResp.Error)
		// c.ctl.Shutdown(emsg)
		fmt.Println(emsg)
		return
	}

	c.id = authResp.ClientId
	c.serverVersion = authResp.MmVersion
	c.Info("Authenticated with server, client id: %v", c.id)
	c.update()
	if err = SaveAuthToken(c.configPath, c.authToken); err != nil {
		c.Error("Failed to save auth token: %v", err)
	}

	// request tunnels
	reqIDToTunnelConfig := make(map[string]*TunnelConfiguration)
	for _, config := range c.tunnelConfig {
		// create the protocol list to ask for
		var protocols []string
		for proto := range config.Protocols {
			protocols = append(protocols, proto)
		}

		reqTunnel := &msg.ReqTunnel{
			ReqId:      util.RandId(8),
			Protocol:   strings.Join(protocols, "+"),
			Hostname:   config.Hostname,
			Subdomain:  config.Subdomain,
			HTTPAuth:   config.HTTPAuth,
			RemotePort: config.RemotePort,
		}

		// send the tunnel request
		if err = msg.WriteMsg(ctlConn, reqTunnel); err != nil {
			panic(err)
		}

		// save request id association so we know which local address
		// to proxy to later
		reqIDToTunnelConfig[reqTunnel.ReqId] = config
	}

	// start the heartbeat
	lastPong := time.Now().UnixNano()
	// c.ctl.Go(func() { c.heartbeat(&lastPong, ctlConn) })
	go func() {
		c.heartbeat(&lastPong, ctlConn)
	}()

	// main control loop
	for {
		var rawMsg msg.Message
		if rawMsg, err = msg.ReadMsg(ctlConn); err != nil {
			panic(err)
		}

		switch m := rawMsg.(type) {
		case *msg.ReqProxy:
			// c.ctl.Go(c.proxy)
			go c.proxy()

		case *msg.Pong:
			atomic.StoreInt64(&lastPong, time.Now().UnixNano())

		case *msg.NewTunnel:
			if m.Error != "" {
				emsg := fmt.Sprintf("Server failed to allocate tunnel: %s", m.Error)
				c.Error(emsg)
				// c.ctl.Shutdown(emsg)
				continue
			}

			tunnel := Tunnel{
				PublicUrl: m.Url,
				LocalAddr: reqIDToTunnelConfig[m.ReqId].Protocols[m.Protocol],
				// Protocol:  c.protoMap[m.Protocol],
				Type: m.Protocol,
			}

			c.tunnels[tunnel.PublicUrl] = tunnel
			c.connStatus = ConnOnline
			c.Info("Tunnel established at %v", tunnel.PublicUrl)
			c.update()

		default:
			log.Warn("Ignoring unknown control message %v ", m)
		}
	}
}

// Establishes and manages a tunnel proxy connection with the server
func (c *Model) proxy() {
	var (
		remoteConn net.Conn
		err        error
	)

	if c.proxyURL == "" {
		remoteConn, err = conn.Dial(c.serverAddr, "pxy")
	} else {
		remoteConn, err = conn.DialHttpProxy(c.proxyURL, c.serverAddr, "pxy")
	}

	if err != nil {
		log.Error("Failed to establish proxy connection: %v", err)
		return
	}
	defer remoteConn.Close()

	err = msg.WriteMsg(remoteConn, &msg.RegProxy{ClientId: c.id})
	if err != nil {
		log.Error("Failed to write RegProxy: %v", err)
		return
	}

	// wait for the server to ack our register
	var startPxy msg.StartProxy
	if err = msg.ReadMsgInto(remoteConn, &startPxy); err != nil {
		log.Error("Server failed to write StartProxy: %v", err)
		return
	}

	tunnel, ok := c.tunnels[startPxy.Url]
	if !ok {
		log.Error("Couldn't find tunnel for proxy: %s", startPxy.Url)
		return
	}

	// start up the private connection
	start := time.Now()
	localConn, err := conn.Dial(tunnel.LocalAddr, "prv")
	if err != nil {
		log.Warn("Failed to open private leg %s: %v", tunnel.LocalAddr, err)

		if tunnel.Type == "http" { //tunnel.Protocol.GetName() == "http"
			// try to be helpful when you're in HTTP mode and a human might see the output
			badGatewayBody := fmt.Sprintf(badGateway, tunnel.PublicUrl, tunnel.LocalAddr, tunnel.LocalAddr)
			remoteConn.Write([]byte(fmt.Sprintf(`HTTP/1.0 502 Bad Gateway
Content-Type: text/html
Content-Length: %d

%s`, len(badGatewayBody), badGatewayBody)))
		}
		return
	}
	defer localConn.Close()

	m := c.metrics
	m.proxySetupTimer.Update(time.Since(start))
	m.connMeter.Mark(1)
	c.update()
	m.connTimer.Time(func() {
		localConn := localConn //tunnel.Protocol.WrapConn(localConn, ConnectionContext{Tunnel: tunnel, ClientAddr: startPxy.ClientAddr})
		bytesIn, bytesOut := conn.Join(localConn, remoteConn)
		m.bytesIn.Update(bytesIn)
		m.bytesOut.Update(bytesOut)
		m.bytesInCount.Inc(bytesIn)
		m.bytesOutCount.Inc(bytesOut)
	})
	c.update()
}

// Hearbeating to ensure our connection ngrokd is still live
func (c *Model) heartbeat(lastPongAddr *int64, conn net.Conn) {
	lastPing := time.Unix(atomic.LoadInt64(lastPongAddr)-1, 0)
	ping := time.NewTicker(pingInterval)
	pongCheck := time.NewTicker(time.Second)

	defer func() {
		ping.Stop()
		pongCheck.Stop()
	}()

	for {
		select {
		case <-pongCheck.C:
			lastPong := time.Unix(0, atomic.LoadInt64(lastPongAddr))
			needPong := lastPong.Sub(lastPing) < 0
			pongLatency := time.Since(lastPing)

			if needPong && pongLatency > maxPongLatency {
				c.Info("Last ping: %v, Last pong: %v", lastPing, lastPong)
				c.Info("Connection stale, haven't gotten PongMsg in %d seconds", int(pongLatency.Seconds()))
				return
			}

		case <-ping.C:
			err := msg.WriteMsg(conn, &msg.Ping{})
			if err != nil {
				log.Debug("Got error %v when writing PingMsg", err)
				return
			}
			lastPing = time.Now()
		}
	}
}
