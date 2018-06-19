package mq

import (
	"context"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	// "github.com/go-redis/redis"
	// "math/rand"
	"time"
)

// type duration struct {
// 	time.Duration
// }
//
// func (d *duration) UnmarshalText(text []byte) error {
// 	if len(text) == 0 {
// 		d.Duration = time.Duration(0)
// 		return nil
// 	}
// 	i, err := strconv.Atoi(string(text))
// 	if err != nil {
// 		return err
// 	}
// 	d.Duration = time.Duration(i) * time.Millisecond
// 	return nil
// }

type RedisConfig struct {
	EtcdEndpoints []string `toml:"etcd_endpoints"`
	EtcdDir       string   `toml:"etcd_dir"`
	// EtcdTTL       duration `toml:"etcd_ttl_ms"`
	EtcdTTLStr string `toml:"etcd_ttl_ms"`
	EtcdTTL    time.Duration
	Member     string `toml:"member"`
	QueueKey   string `toml:"queue_key"`
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

func registerConsumerOneShot(redisConf *RedisConfig) error {
	return registerConsumer(redisConf, false)
}

func registerConsumerKeepAlive(redisConf *RedisConfig) error {
	registerConsumer(redisConf, true)
	ticker := time.NewTicker(redisConf.EtcdTTL)
	defer ticker.Stop()
	// poll
	for {
		fmt.Printf("ttl: %+v\n", redisConf.EtcdTTL)
		fmt.Println("waiting for report")
		_ = <-ticker.C
		fmt.Printf("register consumer for ttl %d\n",
			int(redisConf.EtcdTTL))
		registerConsumer(redisConf, true)
	}
	return nil
}

func registerConsumer(redisConf *RedisConfig, ttl bool) error {
	etcdAPI := getEtcdAPI(redisConf.EtcdEndpoints)
	var setOpt *etcd.SetOptions
	if ttl {
		setOpt = &etcd.SetOptions{TTL: redisConf.EtcdTTL}
	}
	_, err := etcdAPI.Set(context.Background(), redisConf.EtcdDir,
		redisConf.Member, setOpt)
	return err
}

func RegisterConsumer(redisConf *RedisConfig) error {
	if redisConf.EtcdEndpoints == nil {
		if redisConf.EtcdDir != "" || redisConf.EtcdTTLStr != "" {
			// TODO Warning: no etcd endpoints configured but got other etcd
			// configurations
		}
		return nil
	}

	if redisConf.EtcdDir == "" {
		panic("etcd configured, but no etcd dir")
	}

	immediate := time.Duration(0) * time.Millisecond
	// if ttl <= 0, one-shot register
	if redisConf.EtcdTTL <= immediate {
		registerConsumerOneShot(redisConf)
	} else {
		registerConsumerKeepAlive(redisConf)
	}
	return nil
}
