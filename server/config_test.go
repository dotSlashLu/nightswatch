package main

import (
	"fmt"
	"github.com/dotSlashLu/nightswatch/server/message_queue"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const cfgFile = "/tmp/nw-test-config"

func write(cfg string) {
	if err := ioutil.WriteFile(cfgFile, []byte(cfg), 0644); err != nil {
		panic("failed to write test config file " + err.Error())
	}
}

func cleanup() {
	os.Remove(cfgFile)
}

func TestParseConfig(t *testing.T) {
	defer cleanup()
	cfg1 := `
		[log]
		level = ""
		dir = "/var/log/nwatch"

		[message_queue]
		type = "redis"

		[message_queue.redis]
		    # optional etcd configs
		    # if configured, the consumer will register it's consuming queue to the
		    # etcd cluster to let agents know where to push to.
		    # Agents should also be configured to read from this etcd cluster for
		    # target queues.
		    # if unconfigured, agents should be configured directly with queue
		    # addresses
		    etcd_endpoints = ["http://localhost:2379"]
		    etcd_dir = "/nw/redis-servers"
		    etcd_ttl_ms = "1000"

		    # this agent will register and consume from this node
		    # if etcd is not used, fill this directly to agents' config
		    addr = "127.0.0.1:6379"
		    queue_key = "qk-etcd"
		    # the interval of the consumer retrives data from the queue. Increasing
		    # this value will help to increase the throughput.
		    # But if your cluster is busy, decrease this value to reduce the latency.
		    consume_interval = 1000
	`
	cfg1Expect := config{
		logConfig{"", "/var/log/nwatch"},
		mq.Config(messageQueueConfig{
			Type: "redis",
			Conf: mq.RedisConfig{
				EtcdEndpoints: []string{"http://localhost:2379"},
				EtcdDir:       "/nw/redis-servers",
				EtcdTTLStr:    "1000",
				Addr:          "127.0.0.1:6379",
				QueueKey:      "qk-etcd",
			},
		}),
	}
	write(cfg1)
	parsed := parseConfig(cfgFile)
	fmt.Printf("parsed 1 %+v\n", parsed)
	assert.Equal(t, parsed.Log.Level, cfg1Expect.Log.Level, "Log.Level")
	assert.Equal(t, parsed.Log.Dir, cfg1Expect.Log.Dir, "Log.Dir")
	assert.Equal(t, parsed.MessageQueue.Type, cfg1Expect.MessageQueue.Type,
		"MessageQueue.Type")
	etcdEndpointsParsed := (parsed.MessageQueue.Conf.(*mq.RedisConfig).
		EtcdEndpoints)
	etcdEndpointsExp := (cfg1Expect.MessageQueue.Conf.(mq.RedisConfig).
		EtcdEndpoints)
	assert.Equal(t, len(etcdEndpointsExp), len(etcdEndpointsParsed),
		"EtcdEndpoints len")
	if len(etcdEndpointsParsed) > 0 {
		assert.Equal(t, etcdEndpointsExp[0], etcdEndpointsParsed[0],
			"EtcdEndpoints[0]")
	}

	cfg2 := `
		[message_queue]
		type = "redis"
	`
	write(cfg2)
	parsed = parseConfig(cfgFile)
	fmt.Printf("parsed 2 %+v\n", parsed)
	fmt.Printf("parsed 2 Conf: %+v\n", parsed.MessageQueue.Conf)
}
