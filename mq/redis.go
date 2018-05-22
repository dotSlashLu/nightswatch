package nwqueue

import (
	"fmt"
	"github.com/go-redis/redis"
)

type redisQueue struct {
	name string
	conn *redis.Client
}

func getClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Can't connect to redis: %v", err))
	}

	return client
}

func redisGetQueue() *redisQueue {
	client := getClient()
	return &redisQueue{"redis", client}
}

func (q *redisQueue) Push(k string, v string) bool {
	// divide into m queues
	q.conn.RPush(k, v)
	return true
}

func (q *redisQueue) Pop(k string) string {
	// divide into m queues
	v, _ := q.conn.LPop(k).Result()
	return v
}
