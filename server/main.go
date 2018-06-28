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
	flags := parseFlags()
	banner()
	cfg = parseConfig(flags.configFile)
	fmt.Printf("read config %+v\n", cfg)
	queue := mq.New(cfg.MessageQueue)
	queue.StartConsume()
}
