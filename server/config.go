package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dotSlashLu/nightswatch/server/message_queue"
	"io/ioutil"
	"strconv"
	"time"
)

type logConfig struct {
	Level string
	Dir   string
}

type messageQueueConfig struct {
	Type  string `toml:"type"`
	Redis toml.Primitive
	// holds what ever message queue conf Type specified parsed from the above
	// type defs
	Conf interface{}
}

type config struct {
	Log          logConfig
	MessageQueue mq.Config `toml:"message_queue"`
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
		// parse redisConf.EtcdTTLStr to redisConf.EtcdTTL as time.Duration
		// not using UnmarshalText because we have to tell from 0 and empty
		if redisConf.EtcdTTLStr != "" {
			i, err := strconv.Atoi(redisConf.EtcdTTLStr)
			if err != nil {
				panic(fmt.Sprintf("can't parse config field etcd_ttl_ms, %s",
					err))
			}
			redisConf.EtcdTTL = time.Duration(i) * time.Millisecond
		}
		cfg.MessageQueue.Conf = redisConf
	default:
		panic("unrecognized message queue type " + cfg.MessageQueue.Type +
			" in configuration")
	}
	return cfg
}
