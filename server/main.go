package main

import (
	"fmt"
	"github.com/dotSlashLu/nightswatch/server/message_queue"
)

func banner() {
	fmt.Println("Night gathers and my watch begins.")
}

var cfg *config

func main() {
	banner()
	cfg = parseConfig("./etc/server.toml")
	fmt.Printf("read config %+v\n", cfg)
	queue := mq.New(cfg.MessageQueue.Conf.(*mq.RedisConfig))
	queue.StartConsume()
}
