package nwqueue

type NwQueue interface {
	Push(string, string) bool
	Pop(string) string
}

func Init(kind string) NwQueue {
	switch kind {
	case "redis":
		return redisGetQueue()
	default:
		return nil
	}
}
