package nwqueue

import (
	"time"
	"fmt"
	"context"
	"github.com/go-redis/redis"
	etcd "github.com/coreos/etcd/client"
	"math/rand"
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
	EtcdDir string `toml:"etcd_dir"`
	Members  []string
	QueueKey string `toml:"queue_key"`
}

func initClientsByMembers(cfg *RedisConfig) redisClients {
	clients := make([]*redis.Client, len(cfg.Members))
	for i := range cfg.Members {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Members[i],
			Password: "",
			DB:       0,
		})

		_, err := client.Ping().Result()
		if err != nil {
			panic(fmt.Sprintf("Can't connect to redis: %v", err))
		}
		clients[i] = client
	}
	return clients
}

func getEtcdAPI(endpoints []string) etcd.KeysAPI {
	etcdCfg := etcd.Config {
		Endpoints: endpoints,
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
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
	resp, err := api.Get(context.Background(), cfg.EtcdDir, nil)
	if err != nil {
		panic("failed to get " + cfg.EtcdDir + " error: " + err.Error())
	}
	fmt.Printf("got dir: %+v\n", resp)
	return nil
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
	return true
}

// func (q *redisQueue) Pop() (string, string) {
// 	v, _ := q.getConn().LPop(q.cfg.QueueKey).Result()
// 	kvp := strings.Split(v, sep)
// 	k, v := kvp[0], kvp[1]
// 	return k, v
// }
