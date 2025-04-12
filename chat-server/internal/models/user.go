package models

import (
	"time"
	
	"chat-server/internal/utils"
)

// User 表示聊天系统中的用户
type User struct {
	ID        string    `json:"id"`        // 用户唯一ID
	Username  string    `json:"username"`  // 用户名
	Avatar    string    `json:"avatar"`    // 头像URL
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	LastSeen  time.Time `json:"lastSeen"`  // 最后在线时间
	Status    string    `json:"status"`    // 状态: online, offline, away
}

// NewUser 创建一个新用户
func NewUser(username, avatar string) *User {
	return &User{
		ID:        utils.GenerateUUID(),
		Username:  username,
		Avatar:    avatar,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		Status:    "online",
	}
}

// UpdateLastSeen 更新用户最后在线时间
func (u *User) UpdateLastSeen() {
	u.LastSeen = time.Now()
}

// SetStatus 设置用户状态
func (u *User) SetStatus(status string) {
	u.Status = status
	u.UpdateLastSeen()
}

// Room 表示一个聊天室
type Room struct {
	ID          string    `json:"id"`          // 房间唯一ID
	Name        string    `json:"name"`        // 房间名称
	Description string    `json:"description"` // 房间描述
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	CreatedBy   string    `json:"createdBy"`   // 创建者ID
	IsPrivate   bool      `json:"isPrivate"`   // 是否为私有房间
}

// NewRoom 创建一个新房间
func NewRoom(name, description, createdBy string, isPrivate bool) *Room {
	return &Room{
		ID:          utils.GenerateUUID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		IsPrivate:   isPrivate,
	}
}

// GenerateUUID 使用工具包中的函数生成UUID
func GenerateUUID() string {
	return utils.GenerateUUID()
} 