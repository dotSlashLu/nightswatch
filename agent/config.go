package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type configMessageQueue struct {
	Type    string   `toml:"type"`
	Members []string `toml:"members"`
}

type configPlugins struct {
	Directory string   `toml:"dir"`
	Names     []string `toml:"names"`
}

type configLog struct {
	Level     string `toml:"level"`
	Directory string `toml:"dir"`
}

type config struct {
	Log          configLog          `toml:"log"`
	MessageQueue configMessageQueue `toml:"message_queue"`
	Plugins      configPlugins      `toml:"plugins"`
}

func newConfig() *config {
	c := new(config)
	c.Log = configLog{Directory: "/var/log/nwatch"}
	return c
}

func parseConfig(filename string) *config {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("error reading config file " + filename + ", " + err.Error())
	}
	cfg := newConfig()
	if _, err = toml.Decode(string(fileContent), cfg); err != nil {
		panic("error decoding config file: " + err.Error())
	}
	return cfg
}
