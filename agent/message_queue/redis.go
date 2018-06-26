package nwqueue

import (
	"context"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

// separator between metric and value
const sep = "·=·"

type redisClients []*redis.Client

type redisQueue struct {
	name    string
	cfg     *RedisConfig
	connLen int
	conns   redisClients
}

type RedisConfig struct {
	EtcdEndpoints []string `toml:"etcd_endpoints"`
	EtcdDir       string   `toml:"etcd_dir"`
	Members       []string
	QueueKey      string `toml:"queue_key"`
}

func initRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Can't connect to redis: %v", err))
	}
	return client
}

func initClientsByMembers(cfg *RedisConfig) redisClients {
	if len(cfg.Members) < 1 {
		panic("Message queue type configured to redis but no redis server " +
			"configured.")
	}
	clients := make([]*redis.Client, len(cfg.Members))
	for i := range cfg.Members {
		clients[i] = initRedisClient(cfg.Members[i], "", 0)
	}
	return clients
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

func initClientsByEtcd(cfg *RedisConfig) redisClients {
	api := getEtcdAPI(cfg.EtcdEndpoints)
	// read dir
	resp, err := api.Get(context.Background(), cfg.EtcdDir,
		&etcd.GetOptions{Recursive: true})
	if err != nil {
		panic("failed to get " + cfg.EtcdDir + " from etcd error: " +
			err.Error())
	}
	if len(resp.Node.Nodes) < 1 {
		panic("No redis server found in etcd")
	}
	fmt.Printf("got dir: %+v\n", resp.Node.Nodes)
	clients := make([]*redis.Client, len(resp.Node.Nodes))
	for i := range resp.Node.Nodes {
		clients[i] = initRedisClient(resp.Node.Nodes[i].Value, "", 0)
	}
	return clients
}

// members are prioritized
func getClient(cfg *RedisConfig) redisClients {
	var clients redisClients
	if len(cfg.Members) > 0 {
		clients = initClientsByMembers(cfg)
	} else {
		clients = initClientsByEtcd(cfg)
	}
	return clients
}

func redisGetQueue(cfg *RedisConfig) *redisQueue {
	conns := getClient(cfg)
	return &redisQueue{"redis", cfg, len(conns), conns}
}

func (q *redisQueue) getConn() *redis.Client {
	i := rand.Intn(q.connLen)
	conn := q.conns[i]
	return conn
}

func (q *redisQueue) Push(k string, v string) bool {
	kvp := fmt.Sprintf("%s%s%s", k, sep, v)
	q.getConn().RPush(q.cfg.QueueKey, kvp)
	fmt.Printf("push %v to key %v\n", kvp, q.cfg.QueueKey)
	return true
}
