package mq

import (
	"context"
	"fmt"
	"strings"
	etcd "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"time"
)

// separator between metric and value
const sep = "·=·"

type RedisConfig struct {
	EtcdEndpoints []string `toml:"etcd_endpoints"`
	EtcdDir       string   `toml:"etcd_dir"`
	EtcdTTLStr string `toml:"etcd_ttl_ms"`
	EtcdTTL    time.Duration
	Addr       string `toml:"addr"`
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
		redisConf.Addr, setOpt)
	return err
}

func RegisterConsumer(q *redisQueue) error {
	redisConf := q.cfg
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

func initRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Can't connect to redis: %v", err))
	}
	return client
}

type redisQueue struct {
	cfg     *RedisConfig
	conn    *redis.Client
}

func New(config *RedisConfig) *redisQueue {
	q := &redisQueue{ cfg: config }
	client := initRedisClient(config.Addr, "", 0)
	q.conn = client
	return q
}

func (q *redisQueue) StartConsume() {
	RegisterConsumer(q)
}

func (q *redisQueue) Pop() (string, string) {
	v, _ := q.conn.LPop(q.cfg.QueueKey).Result()
	kvp := strings.Split(v, sep)
	k, v := kvp[0], kvp[1]
	return k, v
}

