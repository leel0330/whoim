package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"whoim/common"
)

type ClientConn struct {
	ClientID   int
	ClientAddr string
	Conn       net.Conn
}

type ChatServer struct {
	GateAddr    string
	SrvListener net.Listener

	GroupMap map[int]*ChatGroup

	curClientID int
	curGroupID  int

	sync.RWMutex
}

func NewChatServer(port int) (*ChatServer, error) {
	srv := &ChatServer{
		GateAddr:    fmt.Sprintf(":%v", port),
		curClientID: 0,
		curGroupID:  0,

		GroupMap: make(map[int]*ChatGroup),
	}
	err := srv.Start()
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *ChatServer) Start() error {
	err := srv.startListen()
	if err != nil {
		return err
	}
	return nil
}

func (srv *ChatServer) startListen() error {
	addr, err := net.ResolveTCPAddr("tcp4", srv.GateAddr)
	if err != nil {
		log.Printf("fail to resolve tcp addr:%v", err)
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panicf("fail to init listener")
		return err
	}
	log.Printf("chat service start on addr:%v", addr)
	srv.SrvListener = listener
	return nil
}

func (srv *ChatServer) Run() {
	for {
		conn, err := srv.SrvListener.Accept()
		if err != nil {
			log.Printf("fail to get one conn:%v", err)
			continue
		}
		srv.curClientID += 1
		clientConn := ClientConn{
			ClientID:   srv.curClientID,
			ClientAddr: conn.RemoteAddr().String(),
			Conn:       conn,
		}
		go srv.HandleClientConn(clientConn)
	}
}

func (srv *ChatServer) HandleClientConn(ch ClientConn) {
	log.Printf("new client conn:%v", ch.Conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		size, err := ch.Conn.Read(buf)
		if size == 0 {
			log.Printf("client conn:{%v:%v} closed", ch.ClientID, ch.ClientAddr)
			//如果在一些群组里，需要从群组里删除该对象
			for _, group := range srv.GroupMap {
				srv.LeaveGroup(group.GroupID, ch)
			}
			ch.Conn.Close()
			break
		}
		if err != nil {
			log.Printf("fail to receive data:%v,%v", ch.ClientID, err)
			ch.Conn.Close()
			break
		}

		receiveStr := string(buf[:size])
		srv.handleData(receiveStr, ch)
	}
}

func (srv *ChatServer) handleData(cmd string, ch ClientConn) {
	var msg common.ChatMessage
	err := json.Unmarshal([]byte(cmd), &msg)
	if err != nil {
		log.Printf("fail to unmarshal:%v", cmd)
		return
	}
	log.Printf("msg:%v", msg)
	switch msg.Type {
	case common.CreateGroup:
		srv.curGroupID += 1
		group := NewChatGroup(fmt.Sprintf("g-%v", srv.curGroupID), srv.curGroupID)
		group.AddOneClient(ch)
		srv.GroupMap[group.GroupID] = group
		n, err := ch.Conn.Write([]byte("create group ok"))
		if err != nil {
			log.Printf("fail to create group:%v,%v", n, err)
		}
	case common.AddGroupMember:
		srv.Lock()
		defer srv.Unlock()
		if group, ok := srv.GroupMap[msg.GroupID]; ok {
			group.AddOneClient(ch)
			n, err := ch.Conn.Write([]byte(fmt.Sprintf("add to group %v ok",
				msg.GroupID)))
			if err != nil {
				log.Printf("fail to create group:%v,%v", n, err)
			}
		}
	case common.ShowGroup:
		groups := make([]string, 0)
		for _, group := range srv.GroupMap {
			groups = append(groups, group.GroupName)
		}
		resp := strings.Join(groups, "\n")
		n, err := ch.Conn.Write([]byte(resp))
		if err != nil {
			log.Printf("fail to show group:%v,%v", n, err)
		}
	case common.SHowGroupMembers:
		groupID := msg.GroupID
		members := make([]string, 0)
		if group, ok := srv.GroupMap[groupID]; ok {
			for _, v := range group.ClientMap {
				members = append(members, fmt.Sprintf("%v[%v]", v.ClientID, v.ClientAddr))
			}
		}
		content := strings.Join(members, "\n")
		log.Printf("sgm content:%v", content)
		_, err := ch.Conn.Write([]byte(content))
		if err != nil {
			log.Printf("fail to get group members:%v,%v,%v", groupID, err, members)
		}
	case common.LeaveGroup:
		srv.Lock()
		defer srv.Unlock()
		srv.LeaveGroup(msg.GroupID, ch)
		n, err := ch.Conn.Write([]byte(fmt.Sprintf("you have leave group:%v",
			msg.GroupID)))
		if err != nil {
			log.Printf("fail to leave group:%v,%v", n, err)
		}
	case common.GroupMessage:
		groupID, content := msg.GroupID, msg.Content
		if group, ok := srv.GroupMap[groupID]; ok {
			group.Broadcast(ch, content)
		}
	case common.Normal:
		n, err := ch.Conn.Write([]byte(msg.Content))
		if err != nil {
			log.Printf("fail to send data:%v,%v", n, err)
		}
	}
}

func (srv *ChatServer) LeaveGroup(groupID int, ch ClientConn) {
	if group, ok := srv.GroupMap[groupID]; ok {
		group.RemoveOneClient(ch)
	}
}
