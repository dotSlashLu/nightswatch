// Agent of nightswatch monitoring
package main

import (
	"fmt"
	"github.com/dotSlashLu/nightswatch/common/logger"
	"github.com/dotSlashLu/nightswatch/mq"
	"log"
)

var (
	cfg *config
	q nwqueue.NwQueue
)

func main() {
	flags := parseFlags()
	if !flags.quiet {
		banner()
	}
	cfg = parseConfig(flags.configFile)
	fmt.Printf("read cfg: %+v\n", *cfg)
	logger.Setup("/var/log/nwatch/test.log", 1024*1024)
	q = initQueue()
	loadPlugins()
	fmt.Printf("%v\n", cfg)
}

func init() {
	clientID()
}

func banner() {
	fmt.Println("Night gathers and now my watch begins.")
}

func initQueue() nwqueue.NwQueue {
	mqType := cfg.MessageQueue.Type
	q := nwqueue.Init(mqType)
	log.Printf("inited mq %s\n", mqType)
	return q
}
