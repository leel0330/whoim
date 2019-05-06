package server

import (
	"fmt"
	"log"
	"net"
)

type ClientConn struct {
	ClientID   int
	ClientAddr string
	Conn       net.Conn
}

type GateServer struct {
	GateAddr    string
	SrvListener net.Listener
	curClientID int
	curGroupID  int
}

func NewGateServer(port int) (*GateServer, error) {
	srv := &GateServer{
		GateAddr:    fmt.Sprintf(":%v", port),
		curClientID: 0,
	}
	err := srv.Start()
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *GateServer) Start() error {
	err := srv.startListen()
	if err != nil {
		return err
	}
	return nil
}

func (srv *GateServer) startListen() error {
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
	log.Printf("gate service start on addr:%v", addr)
	srv.SrvListener = listener
	return nil
}

func (srv *GateServer) Run() {
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

func (srv *GateServer) HandleClientConn(ch ClientConn) {
	log.Printf("new client conn:%v", ch.Conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		size, err := ch.Conn.Read(buf)
		if size == 0 {
			log.Printf("client conn:{%v:%v} closed", ch.ClientID, ch.ClientAddr)
			break
		}
		if err != nil {
			log.Printf("fail to receive data:%v,%v", ch.ClientID, err)
			ch.Conn.Close()
			break
		}

		receiveStr := string(buf[:size])
		log.Printf("receive message:%v", receiveStr)

		sendStr := fmt.Sprintf("send to you:%v", receiveStr)
		n, err := ch.Conn.Write([]byte(sendStr))
		if err != nil {
			log.Printf("fail to send msg:%v,%v", ch.ClientID, err)
			ch.Conn.Close()
			break
		}
		log.Printf("send %v size bytes", n)
	}
}
