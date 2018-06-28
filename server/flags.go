package main

import (
	"flag"
)

const defaultConfgFile = "/etc/nwatch/server.toml"

type flags struct {
	configFile string
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.configFile, "c", defaultConfgFile,
		"path to config file")
	flag.Parse()
	return &f
}
