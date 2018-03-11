package httpProxy

import (
	"flag"
	"fmt"
	"mgrok/version"
	"os"
)

type Options struct {
	config  string
	command string
}

func parseArgs() (opts *Options) {
	config := flag.String(
		"config",
		"mgrokp.yaml",
		"Path to httpProxy configuration file. (default: mgrokp.yaml)",
	)

	flag.Parse()

	opts = &Options{
		config:  *config,
		command: flag.Arg(0),
	}

	switch opts.command {
	case "version":
		fmt.Println(version.Full())
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	}

	return opts
}
