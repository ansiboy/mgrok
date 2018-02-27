package conn

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"ngrok/log"
	"sync"
)

// type Conn interface {
// 	net.Conn
// log.Logger
// Id() string
// SetType(string)
// CloseRead() error
// }

type loggedConn struct {
	// tcp *net.TCPConn
	net.Conn
	// log.Logger
	// id  int32
	// typ string
}

type Listener struct {
	net.Addr
	Conns chan net.Conn //*loggedConn
}

// func wrapConn(conn net.Conn, typ string) net.Conn {
// 	return conn
// 	// switch c := conn.(type) {
// 	// case *vhost.HTTPConn:
// 	// 	// wrapped := c.Conn.(*loggedConn)
// 	// 	return &loggedConn{conn} //wrapped.Logger,wrapped.tcp, , rand.Int31(), typ
// 	// case *loggedConn:
// 	// 	return c
// 	// case *net.TCPConn:
// 	// 	wrapped := &loggedConn{conn} //log.NewPrefixLogger(),c, , rand.Int31(), typ
// 	// 	// wrapped.AddLogPrefix(wrapped.Id())
// 	// 	return wrapped
// 	// }

// 	// return nil
// }

func Listen(addr, typ string) (l *Listener, err error) {
	// listen for incoming connections
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	l = &Listener{
		Addr:  listener.Addr(),
		Conns: make(chan net.Conn),
	}

	go func() {
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Error("Failed to accept new TCP connection of type %s: %v", typ, err)
				continue
			}

			c := rawConn //wrapConn(rawConn, typ)
			// if tlsCfg != nil {
			// c.Conn = c.Conn // tls.Server(c.Conn, tlsCfg)
			// }

			log.Info("New connection from %v", c.RemoteAddr())
			l.Conns <- c
		}
	}()
	return
}

// func Wrap(conn net.Conn, typ string) net.Conn {
// 	return wrapConn(conn, typ)
// }

func Dial(addr, typ string) (conn net.Conn, err error) {
	var rawConn net.Conn
	if rawConn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	conn = rawConn //wrapConn(rawConn, typ)
	log.Debug("New connection to: %v", rawConn.RemoteAddr())

	// if tlsCfg != nil {
	// 	conn.StartTLS(tlsCfg)
	// }

	return
}

func DialHttpProxy(proxyUrl, addr, typ string) (conn net.Conn, err error) {
	// parse the proxy address
	var parsedUrl *url.URL
	if parsedUrl, err = url.Parse(proxyUrl); err != nil {
		return
	}

	var proxyAuth string
	if parsedUrl.User != nil {
		proxyAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(parsedUrl.User.String()))
	}

	// var proxyTlsConfig *tls.Config
	// switch parsedUrl.Scheme {
	// case "http":
	// 	proxyTlsConfig = nil
	// case "https":
	// 	proxyTlsConfig = new(tls.Config)
	// default:
	// 	err = fmt.Errorf("Proxy URL scheme must be http or https, got: %s", parsedUrl.Scheme)
	// 	return
	// }

	// dial the proxy
	if conn, err = Dial(parsedUrl.Host, typ); err != nil {
		return
	}

	// send an HTTP proxy CONNECT message
	req, err := http.NewRequest("CONNECT", "https://"+addr, nil)
	if err != nil {
		return
	}

	if proxyAuth != "" {
		req.Header.Set("Proxy-Authorization", proxyAuth)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ngrok)")
	req.Write(conn)

	// read the proxy's response
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non-200 response from proxy server: %s", resp.Status)
		return
	}

	// upgrade to TLS
	// conn.StartTLS(tlsCfg)

	return
}

// func (c *loggedConn) StartTLS(tlsCfg *tls.Config) {
// 	c.Conn = c.Conn // tls.Client(c.Conn, tlsCfg)
// }

// func (c *loggedConn) Close() (err error) {
// 	if err := c.Conn.Close(); err == nil {
// 		log.Debug("Closing")
// 	}
// 	return
// }

// func (c *loggedConn) Id() string {
// 	return fmt.Sprintf("%s:%x", c.typ, c.id)
// }

// func (c *loggedConn) SetType(typ string) {
// 	oldId := c.Id()
// 	c.typ = typ
// 	// c.ClearLogPrefixes()
// 	// c.AddLogPrefix(c.Id())
// 	log.Info("Renamed connection %s", oldId)
// }

// func (c *loggedConn) CloseRead() error {
// 	// XXX: use CloseRead() in Conn.Join() and in Control.shutdown() for cleaner
// 	// connection termination. Unfortunately, when I've tried that, I've observed
// 	// failures where the connection was closed *before* flushing its write buffer,
// 	// set with SetLinger() set properly (which it is by default).
// 	return c.tcp.CloseRead()
// }

func Join(c net.Conn, c2 net.Conn) (int64, int64) {
	var wait sync.WaitGroup

	pipe := func(to net.Conn, from net.Conn, bytesCopied *int64) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		var err error
		*bytesCopied, err = io.Copy(to, from)
		if err != nil {
			log.Warn("Copied %d bytes to %s before failing with error %v", *bytesCopied, "another", err)
		} else {
			log.Debug("Copied %d bytes to %s", *bytesCopied, "another")
		}
	}

	wait.Add(2)
	var fromBytes, toBytes int64
	go pipe(c, c2, &fromBytes)
	go pipe(c2, c, &toBytes)
	log.Info("Joined with connection %s", "another")
	wait.Wait()
	return fromBytes, toBytes
}
