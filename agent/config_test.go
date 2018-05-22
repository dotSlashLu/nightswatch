package main

import "fmt"
import "testing"

const CNF_FILE = "/etc/nwatch/client.toml"

func TestParseConfig(t *testing.T) {
	cfg := parseConfig(CNF_FILE)
	fmt.Printf("parsed %+v\n", cfg)
}
