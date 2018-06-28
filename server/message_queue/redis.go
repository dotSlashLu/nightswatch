package mq

import (
	"context"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

// separator between metric and value
const sep = "·=·"

type consumeInterval struct {
	duration time.Duration
}

type redisQueue struct {
	cfg     *RedisConfig
	conn    *redis.Client
	etcdAPI etcd.KeysAPI
}

func (interval *consumeInterval) UnmarshalText(text []byte) error {
	i, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	interval.duration = time.Duration(i) * time.Millisecond
	return nil
}

type RedisConfig struct {
	EtcdEndpoints   []string `toml:"etcd_endpoints"`
	EtcdDir         string   `toml:"etcd_dir"`
	EtcdTTLStr      string   `toml:"etcd_ttl_ms"`
	EtcdTTL         time.Duration
	Addr            string           `toml:"addr"`
	QueueKey        string           `toml:"queue_key"`
	ConsumeInterval *consumeInterval `toml:"consume_interval"`
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

func registerConsumerOneShot(q *redisQueue) error {
	return registerConsumer(q, false)
}

func registerConsumerKeepAlive(q *redisQueue, firstReg chan bool) error {
	registerConsumer(q, true)
	firstReg <- true
	ticker := time.NewTicker(q.cfg.EtcdTTL)
	defer ticker.Stop()
	// poll
	for {
		_ = <-ticker.C
		fmt.Printf("register consumer for ttl %s\n",
			q.cfg.EtcdTTL)
		registerConsumer(q, true)
	}
	return nil
}

func registerConsumer(q *redisQueue, ttl bool) error {
	etcdAPI := q.etcdAPI
	redisConf := q.cfg
	var setOpt *etcd.SetOptions
	if ttl {
		setOpt = &etcd.SetOptions{TTL: redisConf.EtcdTTL}
	}
	//_, err := etcdAPI.Set(context.Background(), redisConf.EtcdDir,
	//	redisConf.Addr, setOpt)
	key := redisConf.EtcdDir + "/" + redisConf.Addr
	_, err := etcdAPI.Set(context.Background(), key,
		"", setOpt)
	return err
}

func etcdRegister(q *redisQueue, firstReg chan bool) error {
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
		registerConsumerOneShot(q)
		firstReg <- true
	} else {
		registerConsumerKeepAlive(q, firstReg)
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

func newRedis(config *RedisConfig) *redisQueue {
	q := &redisQueue{cfg: config}
	client := initRedisClient(config.Addr, "", 0)
	q.conn = client
	q.etcdAPI = getEtcdAPI(config.EtcdEndpoints)
	return q
}

func (q *redisQueue) StartConsume() {
	firstReg := make(chan bool)
	go etcdRegister(q, firstReg)
	<-firstReg
	close(firstReg)
	fmt.Println("first reg done")
	fmt.Printf("cosume interval: %+v\n", q.cfg.ConsumeInterval.duration)
	ticker := time.NewTicker(q.cfg.ConsumeInterval.duration)
	defer ticker.Stop()
	for {
		<-ticker.C
		records := q.pop()
		for _, r := range records {
			kvp := strings.Split(r, sep)
			k, v := kvp[0], kvp[1]
			fmt.Println("processed", k, v)
		}
		q.trim(int64(len(records)))
	}
}

func (q *redisQueue) pop() []string {
	v, err := q.conn.LRange(q.cfg.QueueKey, 0, -1).Result()
	if err != nil {
		// this just means no value to pop
		if err.Error() == "redis: nil" {
			return []string{}
		}
		panic(err)
	}
	fmt.Printf("redis result: %+v", v)
	return v
}

// trims pop-ed values after processing
func (q *redisQueue) trim(n int64) {
	res, err := q.conn.LTrim(q.cfg.QueueKey, n, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("trimed", res, err)
}
