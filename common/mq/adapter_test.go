package nwqueue

import (
	"fmt"
	"testing"
)

func TestInit(*testing.T) {
	queue := Init("redis")
	fmt.Println(queue)
	fmt.Println("push a, b")
	queue.Push("a", "b")
	fmt.Println("pop a", queue.Pop("a"))
}
