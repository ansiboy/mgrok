package server

import (
	"flag"
)

type Options struct {
	config string
}

func parseArgs() (opts *Options) {
	config := flag.String(
		"config",
		"mgrokd.yaml",
		"Path to mgrok configuration file. (default: mgrok.yaml)",
	)

	opts = &Options{
		config: *config,
	}

	return opts
}
