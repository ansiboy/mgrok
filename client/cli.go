package client

import (
	"flag"
	"fmt"
	"mgrok/version"
	"os"
)

const usage1 string = `Usage: %s [OPTIONS] <local port or address>
Options:
`

const usage2 string = `
Examples:
	mgrok 80
	mgrok -subdomain=example 8080
	mgrok -proto=tcp 22
	mgrok -hostname="example.com" -httpauth="user:password" 10.0.0.1


Advanced usage: mgrok [OPTIONS] <command> [command args] [...]
Commands:
	mgrok start [tunnel] [...]    Start tunnels by name from config file
	mgrok start-all               Start all tunnels defined in config file
	mgrok list                    List tunnel names from config file
	mgrok help                    Print help
	mgrok version                 Print mgrok version

Examples:
	mgrok start www api blog pubsub
	mgrok -log=stdout -config=mgrok.yaml start ssh
	mgrok start-all
	mgrok version

`

// Options options
type Options struct {
	config    string
	log       string
	loglevel  string
	authtoken string
	httpauth  string
	hostname  string
	protocol  string
	subdomain string
	command   string
	args      []string
}

// ParseArgs parse args
func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}

	config := flag.String(
		"config",
		"mgrok.yaml",
		"Path to mgrok configuration file. (default: mgrok.yaml)")

	logto := flag.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings")

	loglevel := flag.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")

	authtoken := flag.String(
		"authtoken",
		"",
		"Authentication token for identifying an mgrok.cn account")

	httpauth := flag.String(
		"httpauth",
		"",
		"username:password HTTP basic auth creds protecting the public tunnel endpoint")

	subdomain := flag.String(
		"subdomain",
		"",
		"Request a custom subdomain from the mgrok server. (HTTP only)")

	hostname := flag.String(
		"hostname",
		"",
		"Request a custom hostname from the mgrok server. (HTTP only) (requires CNAME of your DNS)")

	protocol := flag.String(
		"proto",
		"http",
		"The protocol of the traffic over the tunnel {'http', 'tcp'} (default: 'http')")

	flag.Parse()

	opts = &Options{
		config:    *config,
		log:       *logto,
		loglevel:  *loglevel,
		httpauth:  *httpauth,
		subdomain: *subdomain,
		protocol:  *protocol,
		authtoken: *authtoken,
		hostname:  *hostname,
		command:   flag.Arg(0),
	}

	switch opts.command {
	case "list":
		opts.args = flag.Args()[1:]
	case "start":
		opts.args = flag.Args()[1:]
	case "start-all":
		opts.args = flag.Args()[1:]
	case "version":
		fmt.Println(version.MajorMinor())
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	case "":
		err = fmt.Errorf("Error: Specify a local port to tunnel to, or " +
			"an mgrok command.\n\nExample: To expose port 80, run " +
			"'mgrok 80'")
		return

	default:
		if len(flag.Args()) > 1 {
			err = fmt.Errorf("You may only specify one port to tunnel to on the command line, got %d: %v",
				len(flag.Args()),
				flag.Args())
			return
		}

		opts.command = "default"
		opts.args = flag.Args()
	}

	return
}
