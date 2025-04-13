package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"chat-server/internal/models"
	"chat-server/internal/redis"
)

// Client 表示一个WebSocket客户端连接
type Client struct {
	Hub    *Hub               // 所属的Hub
	Conn   *websocket.Conn    // WebSocket连接
	Send   chan []byte        // 发送消息的通道
	User   *models.User       // 关联的用户信息
	RoomID string             // 当前所在的房间ID
	mu     sync.Mutex         // 保护并发写入
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 取消函数
}

// NewClient 创建一个新的客户端
func NewClient(hub *Hub, conn *websocket.Conn, user *models.User, roomID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256), // 缓冲区大小
		User:   user,
		RoomID: roomID,
		ctx:    ctx,
		cancel: cancel,
	}
}

// ReadPump 持续从WebSocket读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.cancel()
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1024 * 10) // 10KB 消息大小限制
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("websocket read error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		var msg = models.NewDefaultMessage()
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		// 设置发送者ID和房间ID
		msg.SenderID = c.User.ID
		msg.RoomID = c.RoomID

		// 处理特殊命令消息
		if msg.Type == models.CommandMessage {
			c.handleCommand(msg)
			continue
		}

		// 正常消息，广播给所有人
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			continue
		}

		// 保存消息到Redis
		err = redis.SaveMessage(c.ctx, c.RoomID, msgJSON)
		if err != nil {
			log.Printf("failed to save message: %v", err)
		}

		// 发送给Hub广播
		c.Hub.Broadcast <- &BroadcastMessage{
			RoomID:  c.RoomID,
			Message: msgJSON,
		}
	}
}

// WritePump 持续向WebSocket写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中所有消息一次性发送
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 发送Ping保持连接
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// ChangeRoom 切换房间
func (c *Client) ChangeRoom(newRoomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.RoomID == newRoomID {
		return
	}

	// 从旧房间移除
	oldRoomID := c.RoomID
	c.Hub.RemoveFromRoom(c, oldRoomID)

	// 加入新房间
	c.RoomID = newRoomID
	c.Hub.AddToRoom(c, newRoomID)

	// 发送离开消息到旧房间
	leaveMsg := models.NewUserLeaveMessage(oldRoomID, c.User.ID, c.User.Username)
	leaveMsgJSON, _ := json.Marshal(leaveMsg)
	c.Hub.Broadcast <- &BroadcastMessage{
		RoomID:  oldRoomID,
		Message: leaveMsgJSON,
	}

	// 发送加入消息到新房间
	joinMsg := models.NewUserJoinMessage(newRoomID, c.User.ID, c.User.Username)
	joinMsgJSON, _ := json.Marshal(joinMsg)
	c.Hub.Broadcast <- &BroadcastMessage{
		RoomID:  newRoomID,
		Message: joinMsgJSON,
	}
}

// handleCommand 处理命令消息
func (c *Client) handleCommand(msg *models.Message) {
	// 根据命令内容处理不同的指令
	switch msg.Content {
	case "/users":
		// 获取房间用户列表并发送给请求方
		c.sendRoomUsers()
	default:
		// 未知命令，发送错误消息
		errorMsg := models.NewSystemMessage(c.RoomID, "未知命令: "+msg.Content)
		errorMsgJSON, _ := json.Marshal(errorMsg)
		c.Send <- errorMsgJSON
	}
}

// sendRoomUsers 发送房间用户列表
func (c *Client) sendRoomUsers() {
	users := c.Hub.GetRoomUsers(c.RoomID)

	// 创建用户列表消息
	msg := models.NewSystemMessage(c.RoomID, "当前房间用户列表")
	msg.Metadata = users

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal users list: %v", err)
		return
	}

	// 只发送给请求方
	c.Send <- msgJSON
}
