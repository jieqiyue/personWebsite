package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"chat-server/internal/config"
	"chat-server/internal/redis"
	"chat-server/internal/utils"
	"chat-server/internal/websocket"
	"context"
	"log"
	"net/http"
	"strconv"
)

// SetupRouter 初始化Gin路由
func SetupRouter() *gin.Engine {
	// 根据配置设置模式
	if config.AppConfig.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册API路由
	api := router.Group("/api")
	{
		// 健康检查
		api.GET("/ping", Ping)

		// 房间相关API
		rooms := api.Group("/rooms")
		{
			rooms.GET("", GetRooms)
			rooms.POST("", CreateRoom)
			rooms.GET("/:id", GetRoom)
			rooms.GET("/:id/messages", GetRoomMessages)
		}
	}

	// 注册WebSocket路由
	websocket.RegisterRoutes(router)

	return router
}

// Ping 健康检查接口
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"status":  "ok",
	})
}

// GetRooms 获取所有聊天室
func GetRooms(c *gin.Context) {
	rooms := make([]gin.H, 0, len(config.AppConfig.ChatRooms))

	// 从配置中获取聊天室信息
	for _, room := range config.AppConfig.ChatRooms {
		rooms = append(rooms, gin.H{
			"id":          room.ID,
			"name":        room.Name,
			"description": room.Description,
			"isDefault":   room.IsDefault,
			"userCount":   websocket.GlobalHub.GetRoomClientsCount(room.ID),
		})
	}

	c.JSON(200, gin.H{
		"rooms": rooms,
	})
}

// CreateRoom 创建新聊天室
func CreateRoom(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"isPrivate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 生成房间ID
	roomID := utils.GenerateRandomString(8)

	c.JSON(201, gin.H{
		"id":          roomID,
		"name":        req.Name,
		"description": req.Description,
		"isPrivate":   req.IsPrivate,
		"createdAt":   "2023-01-01T00:00:00Z",
	})
}

// GetRoom 获取特定聊天室信息
func GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少房间ID",
		})
		return
	}

	c.JSON(200, gin.H{
		"id":          roomID,
		"name":        "聊天室 " + roomID,
		"description": "聊天室描述",
		"userCount":   websocket.GlobalHub.GetRoomClientsCount(roomID),
		"users":       websocket.GlobalHub.GetRoomUsers(roomID),
	})
}

// GetRoomMessages 获取聊天室消息历史
func GetRoomMessages(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少房间ID",
		})
		return
	}

	// 从查询参数获取limit和offset
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	// 从Redis获取消息历史
	ctx := context.Background()
	messages, err := redis.GetMessages(ctx, roomID, offset, offset+limit-1)
	if err != nil {
		log.Printf("Failed to get messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取消息失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"roomId":   roomID,
		"messages": messages,
		"total":    len(messages),
		"offset":   offset,
		"limit":    limit,
	})
}
