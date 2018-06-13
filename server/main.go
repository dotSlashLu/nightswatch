package main

import (
	"fmt"
)

func banner() {
	fmt.Println("Night gathers and my watch begins.")
}

func main() {
	banner()
	fmt.Printf("read config %+v\n", parseConfig("./etc/server.toml"))
}
