// Agent of nightswatch monitor
// Receivs raven(agent plugin) reports then push then to message queue
package main

import (
	"fmt"
	"github.com/dotSlashLu/nightswatch/agent/message_queue"
	"github.com/dotSlashLu/nightswatch/agent/util"
	"github.com/dotSlashLu/nightswatch/common/logger"
	"log"
	"path/filepath"
)

var (
	cfg *config
	q   nwqueue.NwQueue
)

func main() {
	flags := parseFlags()
	if !flags.quiet {
		banner()
	}
	cfg = parseConfig(flags.configFile)
	fmt.Printf("read cfg: %+v\n", *cfg)
	logger.Setup(cfg.Log.Directory+"/test.log", 1024*1024)
	util.GenerateClientID(filepath.Dir(flags.configFile))
	q = initQueue()
	loadPlugins()
	fmt.Printf("%v\n", cfg)
}

func banner() {
	fmt.Println("Night gathers and now my watch begins.")
}

func initQueue() nwqueue.NwQueue {
	mqType := cfg.MessageQueue.Type
	q := nwqueue.Init(mqType, cfg.MessageQueue.Conf)
	log.Printf("inited mq %s\n", mqType)
	return q
}
