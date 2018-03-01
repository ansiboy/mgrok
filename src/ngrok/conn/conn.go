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

func Dial(addr, typ string) (conn net.Conn, err error) {
	var rawConn net.Conn
	if rawConn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	conn = rawConn
	log.Debug("New connection to: %v", rawConn.RemoteAddr())

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

	return
}

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
