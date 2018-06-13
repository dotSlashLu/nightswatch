package main

import (
	"github.com/BurntSushi/toml"
	nwqueue "github.com/dotSlashLu/nightswatch/agent/message_queue"
	"io/ioutil"
)

type configMessageQueue struct {
	Type string `toml:"type"`

	Redis toml.Primitive
	Conf  interface{}
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
	md, err := toml.Decode(string(fileContent), cfg)
	if err != nil {
		panic("error decoding config file: " + err.Error())
	}
	switch cfg.MessageQueue.Type {
	case "redis":
		redisConf := new(nwqueue.RedisConfig)
		err := md.PrimitiveDecode(cfg.MessageQueue.Redis, redisConf)
		if err != nil {
			panic("can't parse message_queue.redis: " + err.Error())
		}
		cfg.MessageQueue.Conf = redisConf
	}
	return cfg
}
