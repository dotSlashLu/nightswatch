package logger

import (
    "log"
    "testing"
)

func TestSetup(t *testing.T) {
    Setup("./test.log", 10 * 1024 * 1024)
    log.Println("asdfasdfasdf")
    log.Println("asdfasdfasdf")
    return
}
