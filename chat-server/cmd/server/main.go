package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat-server/internal/api"
	"chat-server/internal/config"
	 _ "chat-server/internal/fileprocessor"
	"chat-server/internal/redis"
	"chat-server/internal/websocket"
)

func main() {
	// 加载配置文件
	configPath := config.GetDefaultConfigPath()
	err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("配置加载成功, Debug模式: %v", config.AppConfig.Server.Debug)

	// 初始化Redis连接
	err = redis.Setup()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("Redis连接成功")

	// 处理草稿文件夹中的Markdown文件
	// if err := fileprocessor.ProcessDrafts(); err != nil {
	// 	log.Printf("处理草稿文件失败: %v", err)
	// 	// 继续执行，不因为草稿处理失败而中断整个服务启动
	// }

	// 初始化WebSocket Hub
	websocket.Setup()
	log.Println("WebSocket Hub初始化成功")

	// 设置路由
	router := api.SetupRouter()
	log.Println("API路由初始化成功")

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port),
		Handler: router,
	}

	// 优雅关闭服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 http://%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待终止信号
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭服务器，等待5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}
