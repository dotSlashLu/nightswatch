package main

import (
	"flag"
)

const defaultFile = "/etc/nwatch/server.toml"

type flags struct {
	configFile string
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.configFile, "c", defaultFile,
		"path to config file")
	flag.Parse()
	return &f
}
