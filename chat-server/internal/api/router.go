package api

import (
	"chat-server/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"

	"chat-server/internal/config"
	"chat-server/internal/redis"
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

	router.GET("/version", func(c *gin.Context) {
		data, err := os.ReadFile("/tmp/personWebsite/chat-server/cmd/server/version.txt")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "无法读取版本信息",
				"details": err.Error(),
			})
			return
		}
		c.String(200, string(data))
	})

	// 注册API路由
	api := router.Group("/api")
	{
		// 健康检查
		api.GET("/ping", Ping)

		// 房间相关API
		rooms := api.Group("/rooms")
		{
			rooms.GET("", GetRooms)
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

// GetRoom 获取特定聊天室信息
func GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少房间ID",
		})
		return
	}

	// 检查是否是预配置的普通聊天室
	var foundRoom *config.ChatRoomConfig
	for _, room := range config.AppConfig.ChatRooms {
		if room.ID == roomID {
			foundRoom = &room
			break
		}
	}

	// 如果找到了预配置的普通聊天室，返回其信息
	// 否则返回错误
	if foundRoom != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":          foundRoom.ID,
			"name":        foundRoom.Name,
			"description": foundRoom.Description,
			"isDefault":   foundRoom.IsDefault,
		})
	} else {
		// 如果不是预配置的普通聊天室，返回错误
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在或无权访问",
		})
	}
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
	offsetStr := c.DefaultQuery("offset", "1")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	// 从Redis获取消息历史
	// 首先检查列表长度
	ctx := context.Background()
	msgLen := redis.GetDefaultMessageLen(ctx, roomID)
	var messages []models.Message
	if msgLen > 0 {
		limit = -(limit + offset - 1)
		offset = -offset
		messages, err = redis.GetMessages(ctx, roomID, limit, offset)
		if err != nil {
			log.Printf("Failed to get messages: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取消息失败",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"roomId":   roomID,
		"messages": messages,
		"total":    len(messages),
		"offset":   offset,
		"limit":    limit,
	})
}
