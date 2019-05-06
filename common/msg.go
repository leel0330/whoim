package common

type MessageType int

const (
	Normal MessageType = iota
	CreateGroup
	AddGroupMember
	ShowGroup
	SHowGroupMembers
	LeaveGroup
	GroupMessage
	ShowOnlineUsers
)

type ChatMessage struct {
	Type    MessageType
	GroupID int
	Content string
}
