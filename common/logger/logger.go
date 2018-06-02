package logger

import (
	"log"
)

var Logger *log.Logger

func Setup(filename string, sizeLimit int) {
	Writer, err := NewRotateWriter(filename, sizeLimit)
	if err != nil {
		panic("Error setting up logger: " + err.Error())
	}
	log.SetOutput(Writer)
}
