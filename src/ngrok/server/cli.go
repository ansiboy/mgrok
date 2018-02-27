package server

import (
	"flag"
)

type Options struct {
	HttpAddr         string
	HttpsAddr        string
	TunnelAddr       string
	Domain           string
	TlsCrt           string
	TlsKey           string
	LogTo            string
	LogLevel         string
	HttpPulbishPort  string
	HttpsPulbishPort string
}

// 解释参数
func ParseArgs() *Options {
	httpAddr := flag.String("httpAddr", ":80", "Public address for HTTP connections, empty string to disable")
	httpsAddr := flag.String("httpsAddr", ":443", "Public address listening for HTTPS connections, emptry string to disable")
	tunnelAddr := flag.String("tunnelAddr", ":4443", "Public address listening for ngrok client")
	domain := flag.String("domain", "ngrok.com", "Domain where the tunnels are hosted")
	tlsCrt := flag.String("tlsCrt", "", "Path to a TLS certificate file")
	tlsKey := flag.String("tlsKey", "", "Path to a TLS key file")
	logto := flag.String("log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	loglevel := flag.String("log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	httpPulbishPort := flag.String("httpPulbishPort", "", "Public http port")
	httpsPulbishPort := flag.String("httpsPulbishPort", "", "Public https port")
	flag.Parse()

	return &Options{
		HttpAddr:         *httpAddr,
		HttpsAddr:        *httpsAddr,
		TunnelAddr:       *tunnelAddr,
		Domain:           *domain,
		TlsCrt:           *tlsCrt,
		TlsKey:           *tlsKey,
		LogTo:            *logto,
		LogLevel:         *loglevel,
		HttpPulbishPort:  *httpPulbishPort,
		HttpsPulbishPort: *httpsPulbishPort,
	}
}
