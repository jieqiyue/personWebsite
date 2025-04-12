package models

import (
	"time"
)

// MessageType 消息类型枚举
type MessageType string

const (
	TextMessage    MessageType = "text"     // 文本消息
	ImageMessage   MessageType = "image"    // 图片消息
	SystemMessage  MessageType = "system"   // 系统消息
	JoinMessage    MessageType = "join"     // 加入消息
	LeaveMessage   MessageType = "leave"    // 离开消息
	CommandMessage MessageType = "command"  // 命令消息
)

// Message 表示一条聊天消息
type Message struct {
	ID        string      `json:"id"`         // 消息唯一ID
	RoomID    string      `json:"roomId"`     // 房间ID
	SenderID  string      `json:"senderId"`   // 发送者ID
	Content   string      `json:"content"`    // 消息内容
	Type      MessageType `json:"type"`       // 消息类型
	Timestamp time.Time   `json:"timestamp"`  // 消息时间戳
	Metadata  interface{} `json:"metadata,omitempty"` // 可选的额外数据
}

// NewMessage 创建一条新的消息
func NewMessage(roomID, senderID, content string, msgType MessageType) *Message {
	return &Message{
		ID:        GenerateUUID(),
		RoomID:    roomID,
		SenderID:  senderID,
		Content:   content,
		Type:      msgType,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage 创建一条系统消息
func NewSystemMessage(roomID, content string) *Message {
	return &Message{
		ID:        GenerateUUID(),
		RoomID:    roomID,
		SenderID:  "system",
		Content:   content,
		Type:      SystemMessage,
		Timestamp: time.Now(),
	}
}

// NewUserJoinMessage 创建用户加入消息
func NewUserJoinMessage(roomID, userID, username string) *Message {
	return &Message{
		ID:        GenerateUUID(),
		RoomID:    roomID,
		SenderID:  "system",
		Content:   username + "加入了聊天室",
		Type:      JoinMessage,
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"userId":   userID,
			"username": username,
		},
	}
}

// NewUserLeaveMessage 创建用户离开消息
func NewUserLeaveMessage(roomID, userID, username string) *Message {
	return &Message{
		ID:        GenerateUUID(),
		RoomID:    roomID,
		SenderID:  "system",
		Content:   username + "离开了聊天室",
		Type:      LeaveMessage,
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"userId":   userID,
			"username": username,
		},
	}
} 