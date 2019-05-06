package server

import "sync"

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
	group.ClientMap[clientConn.ClientID] = clientConn
}

func (group *ChatGroup) RemoveOneClient(clientConn ClientConn) {
	group.Lock()
	defer group.Unlock()
	delete(group.ClientMap, clientConn.ClientID)
}

func (group *ChatGroup) Broadcast(clientConn ClientConn, msg string) {
	for _, v := range group.ClientMap {
		if v.ClientID != clientConn.ClientID {
			v.Conn.Write([]byte(msg))
		}
	}
}
