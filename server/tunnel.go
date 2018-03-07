package server

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"mgrok/conn"
	"mgrok/log"
	"mgrok/msg"
	"mgrok/util"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var defaultPortMap = map[string]int{
	"http":  80,
	"https": 443,
	"smtp":  25,
}

// Tunnel A control connection, metadata and proxy connections which
//         route public traffic to a firewalled endpoint.
type Tunnel struct {
	// request that opened the tunnel
	req *msg.ReqTunnel

	// time when the tunnel was opened
	start time.Time

	// public url
	url string

	// tcp listener
	listener *net.TCPListener

	// control connection
	ctl *Control

	// logger
	log.Logger

	// closing
	closing int32
}

// Common functionality for registering virtually hosted protocols
func registerVhost(t *Tunnel, protocol string, servingPort int) (err error) {
	vhost := os.Getenv("VHOST")
	if vhost == "" {
		vhost = fmt.Sprintf("%s:%d", config.Domain, servingPort)
	}

	// Canonicalize virtual host by removing default port (e.g. :80 on HTTP)
	defaultPort, ok := defaultPortMap[protocol]
	if !ok {
		return fmt.Errorf("Couldn't find default port for protocol %s", protocol)
	}

	defaultPortSuffix := fmt.Sprintf(":%d", defaultPort)
	if strings.HasSuffix(vhost, defaultPortSuffix) {
		vhost = vhost[0 : len(vhost)-len(defaultPortSuffix)]
	}

	// Canonicalize by always using lower-case
	vhost = strings.ToLower(vhost)

	// Register for specific hostname
	hostname := strings.ToLower(strings.TrimSpace(t.req.Hostname))
	if hostname != "" {
		t.url = fmt.Sprintf("%s://%s", protocol, hostname)
		return tunnelRegistry.register(t.url, t)
	}

	// Register for specific subdomain
	subdomain := strings.ToLower(strings.TrimSpace(t.req.Subdomain))
	if subdomain != "" {
		t.url = fmt.Sprintf("%s://%s.%s", protocol, subdomain, vhost)
		return tunnelRegistry.register(t.url, t)
	}

	// Register for random URL
	t.url, err = tunnelRegistry.registerRepeat(func() string {
		return fmt.Sprintf("%s://%x.%s", protocol, rand.Int31(), vhost)
	}, t)

	return
}

// Create a new tunnel from a registration message received
// on a control channel
func newTCPTunnel(m *msg.ReqTunnel, ctl *Control) (t *Tunnel, err error) {
	t = &Tunnel{
		req:    m,
		start:  time.Now(),
		ctl:    ctl,
		Logger: log.NewPrefixLogger(),
	}

	bindTCP := func(port int) error {
		if t.listener, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: port}); err != nil {
			err = log.Error("Error binding TCP listener: %v", err)
			return err
		}

		// create the url
		addr := t.listener.Addr().(*net.TCPAddr)
		t.url = fmt.Sprintf("tcp://%s:%d", config.Domain, addr.Port)

		// register it
		if err = tunnelRegistry.registerAndCache(t.url, t); err != nil {
			// This should never be possible because the OS will
			// only assign available ports to us.
			t.listener.Close()
			err = fmt.Errorf("TCP listener bound, but failed to register %s", t.url)
			return err
		}

		go t.listenTCP(t.listener)
		return nil
	}

	// use the custom remote port you asked for
	if t.req.RemotePort != 0 {
		bindTCP(int(t.req.RemotePort))
		return
	}

	// try to return to you the same port you had before
	cachedURL := tunnelRegistry.getCachedRegistration(t)
	if cachedURL != "" {
		var port int
		parts := strings.Split(cachedURL, ":")
		portPart := parts[len(parts)-1]
		port, err = strconv.Atoi(portPart)
		if err != nil {
			log.Error("Failed to parse cached url port as integer: %s", portPart)
		} else {
			// we have a valid, cached port, let's try to bind with it
			if bindTCP(port) != nil {
				log.Warn("Failed to get custom port %d: %v, trying a random one", port, err)
			} else {
				// success, we're done
				return
			}
		}
	}

	// Bind for TCP connections
	bindTCP(0)
	return

}

// Create a new tunnel from a registration message received
// on a control channel
func newHTTPTunnel(m *msg.ReqTunnel, ctl *Control, httpAddr net.Addr) (t *Tunnel, err error) {
	t = &Tunnel{
		req:    m,
		start:  time.Now(),
		ctl:    ctl,
		Logger: log.NewPrefixLogger(),
	}

	proto := t.req.Protocol

	// l, ok := listeners[proto]
	if httpAddr == nil {
		err = fmt.Errorf("Not listening for %s connections", proto)
		return
	}

	servingPort := httpAddr.(*net.TCPAddr).Port //l.Addr.(*net.TCPAddr).Port
	if config.HTTPPulbishPort != "" {
		servingPort, err = strconv.Atoi(config.HTTPPulbishPort)
	}

	if err = registerVhost(t, proto, servingPort); err != nil {
		return
	}

	// pre-encode the http basic auth for fast comparisons later
	if m.HTTPAuth != "" {
		m.HTTPAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(m.HTTPAuth))
	}

	t.AddLogPrefix(t.ctl.id)

	t.Info("Registered new tunnel on: %s", t.ctl.id) // t.ctl.conn.Id())

	metrics.openTunnel(t)
	return
}

func (t *Tunnel) shutdown() {
	t.Info("Shutting down")

	// mark that we're shutting down
	atomic.StoreInt32(&t.closing, 1)

	// if we have a public listener (this is a raw TCP tunnel), shut it down
	if t.listener != nil {
		t.listener.Close()
	}

	// remove ourselves from the tunnel registry
	tunnelRegistry.del(t.url)

	// let the control connection know we're shutting down
	// currently, only the control connection shuts down tunnels,
	// so it doesn't need to know about it
	// t.ctl.stoptunnel <- t

	metrics.closeTunnel(t)
}

func (t *Tunnel) id() string {
	return t.url
}

// Listens for new public tcp connections from the internet.
func (t *Tunnel) listenTCP(listener *net.TCPListener) {
	for {
		defer func() {
			if r := recover(); r != nil {
				log.Warn("listenTcp failed with error %v", r)
			}
		}()

		// accept public connections
		tcpConn, err := listener.AcceptTCP()

		if err != nil {
			// not an error, we're shutting down this tunnel
			if atomic.LoadInt32(&t.closing) == 1 {
				return
			}

			t.Error("Failed to accept new TCP connection: %v", err)
			continue
		}

		conn := tcpConn //conn.Wrap(tcpConn, "pub")
		// conn.AddLogPrefix(t.Id())
		log.Info("New connection from %v", conn.RemoteAddr())

		go t.handlePublicConnection(conn)
	}
}

func (t *Tunnel) handlePublicConnection(publicConn net.Conn) {
	defer publicConn.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Warn("HandlePublicConnection failed with error %v", r)
		}
	}()

	startTime := time.Now()
	metrics.openConnection(t, publicConn)

	var proxyConn net.Conn
	var err error
	for i := 0; i < (2 * proxyMaxPoolSize); i++ {
		// get a proxy connection
		if proxyConn, err = t.ctl.GetProxy(); err != nil {
			t.Warn("Failed to get proxy connection: %v", err)
			return
		}
		defer proxyConn.Close()
		t.Info("Got proxy connection %s", "another") //proxyConn.Id())
		// proxyConn.AddLogPrefix(t.Id())

		// tell the client we're going to start using this proxy connection
		startPxyMsg := &msg.StartProxy{
			Url:        t.url,
			ClientAddr: publicConn.RemoteAddr().String(),
		}

		if err = msg.WriteMsg(proxyConn, startPxyMsg); err != nil {
			log.Warn("Failed to write StartProxyMessage: %v, attempt %d", err, i)
			proxyConn.Close()
		} else {
			// success
			break
		}
	}

	if err != nil {
		// give up
		log.Error("Too many failures starting proxy connection")
		return
	}

	// To reduce latency handling tunnel connections, we employ the following curde heuristic:
	// Whenever we take a proxy connection from the pool, replace it with a new one
	util.PanicToError(func() { t.ctl.out <- &msg.ReqProxy{} })

	// no timeouts while connections are joined
	proxyConn.SetDeadline(time.Time{})

	// join the public and proxy connections
	bytesIn, bytesOut := conn.Join(publicConn, proxyConn)
	metrics.closeConnection(t, publicConn, startTime, bytesIn, bytesOut)
}
