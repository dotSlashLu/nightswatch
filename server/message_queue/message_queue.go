package mq

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Type  string `toml:"type"`
	Redis toml.Primitive
	// holds what ever message queue conf Type specified parsed from the above
	// type defs
	Conf interface{}
}

type MessageQueue interface {
	StartConsume()
}

func New(cfg Config) MessageQueue {
	switch cfg.Type {
	case "redis":
		return newRedis(cfg.Conf.(*RedisConfig))
	// this is actually not necessary since queue type has already been
	// verified when parsing config
	// left it here as a defense anyway
	default:
		err := fmt.Sprintf("Message queue configured as %s, "+
			"but not implemented", cfg.Type)
		panic(err)
	}
}
