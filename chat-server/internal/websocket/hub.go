package websocket

import (
	"encoding/json"
	"log"
	"sync"
	
	"chat-server/internal/models"
)

// BroadcastMessage 表示一条广播消息
type BroadcastMessage struct {
	RoomID  string  // 目标房间ID
	Message []byte  // 消息内容
}

// Hub 管理所有WebSocket客户端
type Hub struct {
	// 所有活跃的客户端
	clients map[*Client]bool
	
	// 房间映射，每个房间包含多个客户端
	rooms map[string]map[*Client]bool
	
	// 广播消息通道
	Broadcast chan *BroadcastMessage
	
	// 注册客户端通道
	Register chan *Client
	
	// 取消注册客户端通道
	Unregister chan *Client
	
	// 互斥锁保护map操作
	mu sync.RWMutex
}

// NewHub 创建一个新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan *BroadcastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run 启动Hub处理循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		
		case client := <-h.Unregister:
			h.unregisterClient(client)
		
		case msg := <-h.Broadcast:
			h.broadcastMessage(msg)
		}
	}
}

// registerClient 注册一个新客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// 添加到全局客户端列表
	h.clients[client] = true
	
	// 添加到房间
	h.AddToRoom(client, client.RoomID)
	
	// 广播用户加入消息
	joinMsg := models.NewUserJoinMessage(client.RoomID, client.User.ID, client.User.Username)
	msgJSON, _ := json.Marshal(joinMsg)
	
	// 不需要使用h.Broadcast通道，直接调用广播方法
	h.broadcastToRoom(client.RoomID, msgJSON)
	
	log.Printf("Client registered: %s in room %s", client.User.Username, client.RoomID)
}

// unregisterClient 注销一个客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, ok := h.clients[client]; ok {
		// 从全局客户端列表移除
		delete(h.clients, client)
		
		// 关闭发送通道
		close(client.Send)
		
		// 从房间移除
		h.RemoveFromRoom(client, client.RoomID)
		
		// 广播用户离开消息
		leaveMsg := models.NewUserLeaveMessage(client.RoomID, client.User.ID, client.User.Username)
		msgJSON, _ := json.Marshal(leaveMsg)
		
		// 直接调用广播方法
		h.broadcastToRoom(client.RoomID, msgJSON)
		
		log.Printf("Client unregistered: %s from room %s", client.User.Username, client.RoomID)
	}
}

// broadcastMessage 广播消息到指定房间
func (h *Hub) broadcastMessage(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	h.broadcastToRoom(msg.RoomID, msg.Message)
}

// broadcastToRoom 向房间内所有客户端发送消息
func (h *Hub) broadcastToRoom(roomID string, message []byte) {
	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
				// 消息发送成功
			default:
				// 客户端缓冲区已满，关闭连接
				close(client.Send)
				delete(h.clients, client)
				h.RemoveFromRoom(client, roomID)
			}
		}
	}
}

// AddToRoom 将客户端添加到房间
func (h *Hub) AddToRoom(client *Client, roomID string) {
	// 获取或创建房间
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	
	// 添加客户端到房间
	h.rooms[roomID][client] = true
}

// RemoveFromRoom 将客户端从房间移除
func (h *Hub) RemoveFromRoom(client *Client, roomID string) {
	if _, ok := h.rooms[roomID]; ok {
		delete(h.rooms[roomID], client)
		
		// 如果房间空了，删除房间
		if len(h.rooms[roomID]) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// GetRoomUsers 获取房间内的所有用户
func (h *Hub) GetRoomUsers(roomID string) []*models.User {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	users := make([]*models.User, 0)
	
	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			users = append(users, client.User)
		}
	}
	
	return users
}

// GetRoomClientsCount 获取房间内客户端数量
func (h *Hub) GetRoomClientsCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	
	return 0
}

// RoomExists 检查指定ID的房间是否存在
func (h *Hub) RoomExists(roomID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	_, exists := h.rooms[roomID]
	return exists
} 