package main

import (
	"flag"
)

type flags struct {
	configFile string
	quiet      bool
}

func parseFlags() *flags {
	f := flags{}
	defaultConfigFile := "/etc/nwatch/agent.toml"
	flag.StringVar(&f.configFile, "c", defaultConfigFile, "path to config file")
	flag.BoolVar(&f.quiet, "q", false, "quiet output")
	flag.Parse()
	return &f
}
