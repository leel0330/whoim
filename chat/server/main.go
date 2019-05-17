package main

import (
	"flag"
	"log"
	"whoim/chat/server/services"
)

var (
	port = flag.Int("chat_port", 8000, "chat service port")
)

func main() {

	chatServer, err := services.NewChatServer(*port)
	if err != nil {
		log.Printf("fail to start chat service:%v", err)
		return
	}
	chatServer.Run()
}
