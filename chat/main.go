package main

import (
	"flag"
	"log"
	"whoim/chat/server"
)

var (
	port = flag.Int("chat_port", 8000, "chat service port")
)

func main() {

	chatServer, err := server.NewChatServer(*port)
	if err != nil {
		log.Printf("fail to start chat service:%v", err)
		return
	}
	chatServer.Run()
}
