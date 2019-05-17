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
	ChatTo
)

const (
	SepStr = "\r\n\r\n"
)

type ChatMessage struct {
	Type    MessageType
	GroupID int
	UserId int
	Content string
}
