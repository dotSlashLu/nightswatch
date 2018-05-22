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

type config struct {
	MessageQueue configMessageQueue `toml:"message_queue"`
	Plugins      configPlugins      `toml:"plugins"`
}

func parseConfig(filename string) *config {
	// file, err := os.Open(filename)
	// if err != nil {
	// 	panic("error reading config file " + filename + ", " + err.Error())
	// }
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("error reading config file " + filename + ", " + err.Error())
	}
	cfg := config{}
	if _, err = toml.Decode(string(fileContent), &cfg); err != nil {
		panic("error decoding config file: " + err.Error())
	}
	return &cfg
}
