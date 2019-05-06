package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	host = flag.String("host", "127.0.0.1", "chat service host")
	port = flag.Int("port", 8000, "chat service port")
)

func PrintInnerCommands() {
	cmdStr := `
		:q Quit
		:cg Create Group
		:sg Show Groups
		:su Show Online Users
	`
	log.Printf(cmdStr)
}

func HandleSend(conn *net.TCPConn) {
	var input string
	username := conn.LocalAddr().String()
	for {
		PrintInnerCommands()
		fmt.Scanln(&input)
		if input == ":q" {
			log.Printf("Byebye...")
			conn.Close()
			break
		}
		words := fmt.Sprintf("[%v] say {%v}", username, input)
		size, err := conn.Write([]byte(words))
		if err != nil {
			log.Printf("fail to send msg:%v", err)
			conn.CloseWrite()
			break
		}
		log.Printf("send msg sz:%v,%v", size, words)
	}
}

func main() {
	serverAddrStr := fmt.Sprintf("%v:%v", *host, *port)
	serverAddr, err := net.ResolveTCPAddr("tcp4", serverAddrStr)
	if err != nil {
		log.Printf("fail to resolve tcp addr")
		return
	}
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Printf("fail to connect service:%v", err)
		return
	}

	log.Printf("client conn to service:%v", conn.RemoteAddr().String())

	go HandleSend(conn)

	buf := make([]byte, 1024)
	for {
		size, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		log.Printf("receive data from service:%v", string(buf[:size]))
	}

}
