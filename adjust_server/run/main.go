package main

import (
	"github.com/NeoHuang/adjust-mini/adjust_server"
)

var Version string

func main() {
	adjust_server.New(Version).Start()
}
