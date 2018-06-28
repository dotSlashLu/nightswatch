package main

import (
	"flag"
)

const defaultConfigFile = "/etc/nwatch/agent.toml"

type flags struct {
	configFile string
	quiet      bool
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.configFile, "c", defaultConfigFile,
		"path to config file")
	flag.BoolVar(&f.quiet, "q", false, "quiet output")
	flag.Parse()
	return &f
}
