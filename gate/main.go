package main

import (
	"flag"
	"gochat/gate/server"
	"log"
)

var (
	port = flag.Int("gate_port", 8000, "gate service port")
)

func main() {

	gateServer, err := server.NewGateServer(*port)
	if err != nil {
		log.Printf("fail to start gate service:%v", err)
		return
	}
	gateServer.Run()
}
