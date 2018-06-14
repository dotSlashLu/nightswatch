package mq

import (
	"context"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	d.Duration = time.Duration(nms) * time.Millisecond
	return nil
}

type RedisConfig struct {
	EtcdEndpoints []string `toml:"etcd_endpoints"`
	EtcdDir       string   `toml:"etcd_dir"`
	EtcdTTL       duration `toml:"etcd_ttl_ms"`
	Member        []string `toml:"member"`
	QueueKey      string   `toml:"queue_key"`
}

func getEtcdAPI(endpoints []string) etcd.KeysAPI {
	etcdCfg := etcd.Config{
		Endpoints: endpoints,
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast
		// when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := etcd.New(etcdCfg)
	if err != nil {
		panic("can't connect to etcd error: " + err.Error())
	}
	api := etcd.NewKeysAPI(client)
	return api
}

func registerConsumer() error {
	etcdAPI = getEtcdAPI(redisConf.EtcdEndpoints)
	setOpt = &etcd.SetOptions{TTL: redisConf.EtcdTTL}
	resp, err := etcdAPI.Set(context.Background(), redisConf.EtcdDir,
		redisConf.Member, setOpt)
	return nil
}

func RegisterConsumer(redisConf RedisConfig) error {
	if !redisConf.EtcdEndpoints {
		if redisConf.EtcdDir || redisConf.EtcdTTL {
			// TODO Warning: no etcd endpoints configured but got other etcd
			// configurations
		}
		return nil
	}

	if !redisConf.EtcdDir {
		panic("etcd configured, but no etcd dir")
	}

	// ignore <= 0
	immediate := time.Duration(0) * time.Millisecond
	if redisConf.EtcdTTL > immediate {
	}
}

