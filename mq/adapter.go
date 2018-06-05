// Message Queue Adapter for Nightswatch Monitoring
// To register a new type of message queue:
//	- implement nwqueue.NwQueue interface
//	- add to nwqueue.Init
//	- define nwqueue.QueueNameConfig
//	- register in config parser
package nwqueue

type NwQueue interface {
	Push(string, string) bool
	Pop(string) string
}

func Init(kind string, cfg interface{}) NwQueue {
	switch kind {
	case "redis":
		return redisGetQueue(cfg.(*RedisConfig))
	default:
		return nil
	}
}
