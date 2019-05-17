package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"whoim/chat/common"
)

var (
	host = flag.String("host", "127.0.0.1", "chat service host")
	port = flag.Int("port", 8000, "chat service port")
)

func PrintInnerCommands() {
	cmdStr := `
		:q               Quit
		:cg              Create Group
		:ag	 {gid}		 Add Group
		:sg              Show Groups
		:sgm {gid}       Show Group Members
		:lg {gid}        Leave Group
		:ggm {gid} {msg} Send Group Message
		:sou             Show Online Users
		:ct {uid} {msg}  Chat to user{uid}
	`
	log.Printf(cmdStr)
}

func genChatMessage(cmdStr string) common.ChatMessage {
	tokens := strings.Split(cmdStr, " ")
	log.Printf("tokens:%v", tokens)
	msg := common.ChatMessage{}
	switch tokens[0] {
	case ":cg":
		msg.Type = common.CreateGroup
	case ":ag":
		msg.Type = common.AddGroupMember
		groupID, _ := strconv.Atoi(tokens[1])
		msg.GroupID = groupID
	case ":sg":
		msg.Type = common.ShowGroup
	case ":sgm":
		msg.Type = common.SHowGroupMembers
		groupID, _ := strconv.Atoi(tokens[1])
		msg.GroupID = groupID
	case ":lg":
		msg.Type = common.LeaveGroup
		groupID, _ := strconv.Atoi(tokens[1])
		msg.GroupID = groupID
	case ":ggm":
		msg.Type = common.GroupMessage
		groupID, _ := strconv.Atoi(tokens[1])
		msg.GroupID = groupID
		msg.Content = tokens[2]
	case ":sou":
		msg.Type = common.ShowOnlineUsers
	case ":ct":
		msg.Type = common.ChatTo
		uid, _ := strconv.Atoi(tokens[1])
		msg.UserId = uid
		msg.Content = strings.Join(tokens[2:], " ")
	default:
		msg.Type = common.Normal
		msg.Content = cmdStr
	}
	return msg
}

func SendMessage(conn *net.TCPConn) {
	var input *bufio.Reader
	for {
		PrintInnerCommands()
		input = bufio.NewReader(os.Stdin)
		inputStr, err := input.ReadString('\n')
		if err != nil {
			log.Printf("fail to read stdin input:%v", err)
			continue
		}
		inputStr = inputStr[:len(inputStr)-1]
		if inputStr == ":q" {
			log.Printf("Byebye...")
			err := conn.CloseWrite()
			if err != nil {
				log.Printf("fail to close connection:%v", err)
			}
			break
		}
		msg := genChatMessage(inputStr)
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("fail to json marshal:%v", err)
			continue
		}
		inputStr = string(bytes)
		inputStr = strings.Join([]string{inputStr, common.SepStr}, "")
		log.Printf("input str:%v", inputStr)
		n, err := conn.Write([]byte(inputStr))
		if err != nil {
			log.Printf("fail to send msg:%v,%v", n, err)
			err := conn.CloseWrite()
			if err != nil {
				log.Printf("fail to close connection:%v", err)
			}
			break
		}
	}
}

func genFakeContent() string {
	chars := "abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	charsLen := len(chars)
	s, e := rand.Intn(charsLen), rand.Intn(charsLen)
	if s > e {
		s, e = e, s
	}
	return chars[s : e+1]

}

func AutoSendMessage(conn *net.TCPConn) {
	for {
		ticker := time.NewTicker(time.Second * 1)

		select {
		case _, ok := <-ticker.C:
			if ok {
				content := genFakeContent()
				msg := common.ChatMessage{
					Type:    common.Normal,
					Content: content,
				}
				bytes, err := json.Marshal(msg)
				if err != nil {
					log.Printf("fail to json marshal:%v", err)
					continue
				}
				inputStr := string(bytes)
				inputStr = strings.Join([]string{inputStr, common.SepStr}, "")
				log.Printf("input str:%v", inputStr)
				n, err := conn.Write([]byte(inputStr))
				if err != nil {
					log.Printf("fail to send msg:%v,%v", n, err)
					err := conn.CloseWrite()
					if err != nil {
						log.Printf("fail to close connection:%v", err)
					}
					break
				}

			}


		}
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

	log.Printf("connect to server:%v", conn.RemoteAddr().String())

	go AutoSendMessage(conn)

	buf := make([]byte, 1024)
	for {
		size, err := conn.Read(buf)
		if err != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("fail to close connection:%v", err)
			}
			break
		}
		log.Printf("receive data from server:\n%v", string(buf[:size]))
	}

}
