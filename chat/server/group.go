package server

import (
	"log"
	"sync"
)

type ChatGroup struct {
	GroupID   int
	GroupName string
	ClientMap map[int]ClientConn
	sync.RWMutex
}

func NewChatGroup(name string, groupID int) *ChatGroup {
	group := &ChatGroup{
		GroupName: name,
		GroupID:   groupID,
		ClientMap: make(map[int]ClientConn),
	}
	return group
}

func (group *ChatGroup) AddOneClient(clientConn ClientConn) {
	group.Lock()
	defer group.Unlock()
	log.Printf("client id:%v", clientConn.ClientID)
	group.ClientMap[clientConn.ClientID] = clientConn
}

func (group *ChatGroup) RemoveOneClient(clientConn ClientConn) {
	group.Lock()
	defer group.Unlock()
	if _, ok := group.ClientMap[clientConn.ClientID]; ok {
		delete(group.ClientMap, clientConn.ClientID)
	}
}

func (group *ChatGroup) Broadcast(clientConn ClientConn, msg string) {
	for _, v := range group.ClientMap {
		if v.ClientID != clientConn.ClientID {
			n, err := v.Conn.Write([]byte(msg))
			if err != nil {
				log.Printf("fail to broadcast in group:%v,%v, %v",
					group.GroupID, n, err)
			}
		}
	}
}
