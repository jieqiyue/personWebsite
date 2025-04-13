package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chat-server/internal/config"
	"chat-server/internal/models"
	"chat-server/internal/utils"
)

var (
	// 升级HTTP连接到WebSocket的升级器
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有来源的跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 全局Hub实例
	GlobalHub *Hub
)

// Setup 初始化WebSocket处理
func Setup() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
}

// HandleConnection 处理WebSocket连接请求
func HandleConnection(c *gin.Context) {
	// 简化用户连接方式，用户可以直接提供昵称连接，不需要登录
	// 如果没有提供，则自动生成昵称
	username := c.DefaultQuery("username", "游客"+utils.GenerateRandomString(5))

	// 获取用户请求的房间ID
	roomID := c.DefaultQuery("roomId", "")

	// 验证房间ID是否有效，如果无效则使用默认房间
	if roomID == "" || !config.IsValidRoomID(roomID) {
		roomID = config.GetDefaultRoomID()
		log.Printf("用户请求的房间ID无效，使用默认房间: %s", roomID)
	}

	// 为新用户生成一个临时ID
	userID := utils.GenerateUUID()

	// 创建用户对象
	user := &models.User{
		ID:        userID,
		Username:  username,
		Avatar:    c.DefaultQuery("avatar", ""),
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		Status:    "online",
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// 创建新的客户端
	client := NewClient(GlobalHub, conn, user, roomID)

	// 将客户端注册到Hub
	GlobalHub.Register <- client

	// 启动客户端的读写协程
	go client.ReadPump()
	go client.WritePump()

	log.Printf("New WebSocket connection: user=%s, room=%s", username, roomID)
}

// GetWebSocketInfo 获取WebSocket状态信息
func GetWebSocketInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"activeRooms": len(GlobalHub.rooms),
		"activeUsers": len(GlobalHub.clients),
		"serverTime":  time.Now(),
		"defaultRoom": "general",
	})
}

// RegisterRoutes 注册WebSocket路由
func RegisterRoutes(router *gin.Engine) {
	wsPath := config.AppConfig.WebSocket.Path

	// WebSocket连接处理
	router.GET(wsPath, HandleConnection)

	// WebSocket状态信息API
	router.GET("/api/ws/info", GetWebSocketInfo)
}
