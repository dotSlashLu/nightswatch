package nwqueue

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"strings"
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
	Members  []string
	QueueKey string `toml:"queue_key"`
}

func getClient(cfg *RedisConfig) redisClients {
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

func redisGetQueue(cfg *RedisConfig) *redisQueue {
	conns := getClient(cfg)
	return &redisQueue{"redis", cfg, len(conns), conns}
}

func (q *redisQueue) getConn() *redis.Client {
	i := rand.Intn(q.connLen)
	conn := q.conns[i]
	fmt.Printf("got idx: %d val: %v\n", i, conn)
	return conn
}

func (q *redisQueue) Push(k string, v string) bool {
	kvp := fmt.Sprintf("%s%s%s", k, sep, v)
	q.getConn().RPush(q.cfg.QueueKey, kvp)
	return true
}

func (q *redisQueue) Pop() (string, string) {
	v, _ := q.getConn().LPop(q.cfg.QueueKey).Result()
	kvp := strings.Split(v, sep)
	k, v := kvp[0], kvp[1]
	return k, v
}
