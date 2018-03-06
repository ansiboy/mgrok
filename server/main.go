package server

import (
	"fmt"
	"math/rand"
	log "mgrok/log"
	"mgrok/msg"
	"mgrok/util"
	"net"
	"os"
	"runtime/debug"
	"time"
)

const (
	registryCacheSize uint64        = 1024 * 1024 // 1 MB
	connReadTimeout   time.Duration = 10 * time.Second
)

// GLOBALS
var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	opts *Configuration
	// listeners map[string]*conn.Listener
)

// NewProxy new proxy
func newProxy(pxyConn net.Conn, regPxy *msg.RegProxy) {
	// fail gracefully if the proxy connection fails to register
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Failed with error: %v", r)
			pxyConn.Close()
		}
	}()

	// set logging prefix
	// pxyConn.SetType("pxy")

	// look up the control connection for this proxy
	log.Info("Registering new proxy for %s", regPxy.ClientId)
	ctl := controlRegistry.get(regPxy.ClientId)

	if ctl == nil {
		panic("No client found for identifier: " + regPxy.ClientId)
	}

	ctl.RegisterProxy(pxyConn)
}

// Listen for incoming control and proxy connections
// We listen for incoming control and proxy connections on the same port
// for ease of deployment. The hope is that by running on port 443, using
// TLS and running all connections over the same port, we can bust through
// restrictive firewalls.
func tunnelListener(tunnelAddr string, httpAddr net.Addr) {
	// listen for incoming connections
	listener, err := net.Listen("tcp", tunnelAddr)
	if err != nil {
		panic(err)
	}

	log.Info("Listening for control and proxy connections on %s", tunnelAddr)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Error("Failed to accept new TCP connection of type tcp: %v", err)
			continue
		}

		go tunnelHandler(c, httpAddr)

		log.Info("New connection from %v", c.RemoteAddr())
	}

}

func tunnelHandler(tunnelConn net.Conn, httpAddr net.Addr) {
	// don't crash on panics
	defer func() {
		if r := recover(); r != nil {
			log.Info("tunnelListener failed with error %v: %s", r, debug.Stack())
		}
	}()

	tunnelConn.SetReadDeadline(time.Now().Add(connReadTimeout))
	var rawMsg msg.Message
	var err error
	if rawMsg, err = msg.ReadMsg(tunnelConn); err != nil {
		log.Warn("Failed to read message: %v", err)
		tunnelConn.Close()
		return
	}

	// don't timeout after the initial read, tunnel heartbeating will kill
	// dead connections
	tunnelConn.SetReadDeadline(time.Time{})

	switch m := rawMsg.(type) {
	case *msg.Auth:
		newControl(tunnelConn, m, httpAddr)

	case *msg.RegProxy:
		newProxy(tunnelConn, m)

	default:
		tunnelConn.Close()
	}
}

// Main server main
func Main() {
	// read configuration file
	config, err := LoadConfiguration("")
	opts = config

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// init logging
	log.LogTo(opts.LogTo, opts.LogLevel)

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		panic(err)
	}
	rand.Seed(seed)

	var redirectData *TunnelCenterRegistry
	if config.DataAddr != "" {
		redirectData, err = newRedirectData(config.DataAddr, config.HTTPAddr)
		if err != nil {
			os.Exit(1)
		}
	}

	// init tunnel/control registry
	registryCacheFile := os.Getenv("REGISTRY_CACHE_FILE")
	tunnelRegistry = newTunnelRegistry(registryCacheSize, registryCacheFile, redirectData)
	controlRegistry = newControlRegistry()

	// listen for http
	var httpAddr net.Addr
	if opts.HTTPAddr != "" {
		httpAddr = startHTTPListener(opts.HTTPAddr)
	}

	// ngrok clients
	tunnelListener(opts.TunnelAddr, httpAddr)
	fmt.Scanln()
}
