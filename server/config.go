package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"github.com/dotslashLu/nightswatch/server/message_queue"
)

type logConfig struct {
	Level string
	Dir   string
}

type messageQueueConfig struct {
	Type string `toml:"type"`

	Redis toml.Primitive

	// holds what ever message queue conf Type specified parsed from the above
	// type defs
	Conf interface{}
}

type config struct {
	Log          logConfig
	MessageQueue messageQueueConfig `toml:"message_queue"`
}

func parseConfig(filename string) *config {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("error reading config file " + filename + ", " + err.Error())
	}

	cfg := new(config)
	md, err := toml.Decode(string(fileContent), cfg)
	if err != nil {
		panic("error decoding config file: " + err.Error())
	}
	switch cfg.MessageQueue.Type {
	case "redis":
		redisConf := new(mq.RedisConfig)
		err := md.PrimitiveDecode(cfg.MessageQueue.Redis, redisConf)
		if err != nil {
			panic("can't parse message_queue.redis: " + err.Error())
		}
		cfg.MessageQueue.Conf = redisConf
	default:
		panic("unrecognized message queue type " + cfg.MessageQueue.Type)
	}
	return cfg
}
