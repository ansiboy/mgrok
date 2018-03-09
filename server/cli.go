package server

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
		"mgrokd.yaml",
		"Path to mgrok configuration file. (default: mgrok.yaml)",
	)

	flag.Parse()

	opts = &Options{
		config:  *config,
		command: flag.Arg(0),
	}

	switch opts.command {
	case "version":
		fmt.Println(version.MajorMinor())
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	}

	return opts
}
