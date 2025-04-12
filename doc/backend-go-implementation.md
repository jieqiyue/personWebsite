# Go后端实现设计

本文档详细描述了使用Go语言实现聊天室和在线人数统计功能的后端服务设计。

## 1. 项目结构

```
chat-server/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口点
├── internal/
│   ├── api/
│   │   ├── handlers.go             # HTTP处理器
│   │   └── routes.go               # 路由配置
│   ├── config/
│   │   └── config.go               # 配置管理
│   ├── models/
│   │   ├── message.go              # 消息模型
│   │   └── user.go                 # 用户模型
│   ├── redis/
│   │   └── client.go               # Redis客户端
│   ├── websocket/
│   │   ├── client.go               # WebSocket客户端连接
│   │   ├── hub.go                  # 连接管理中心
│   │   └── message.go              # WebSocket消息处理
│   └── utils/
│       ├── security.go             # 安全工具函数
│       └── id.go                   # ID生成
├── configs/
│   └── config.yaml                 # 配置文件
├── go.mod                          # Go模块文件
└── go.sum                          # 依赖校验文件
```

## 2. 核心组件

### 2.1 WebSocket Hub (连接管理中心)

Hub负责管理所有WebSocket连接，处理用户加入/离开，并广播消息。

```go
// internal/websocket/hub.go

package websocket

import (
    "sync"
    "time"
    
    "chat-server/internal/models"
    "chat-server/internal/redis"
)

// Hub 维护活跃WebSocket连接的集合，广播消息给所有客户端
type Hub struct {
    // 所有活跃的连接
    clients map[*Client]bool

    // 用户数据，以用户ID为键
    users map[string]*models.User

    // 向所有客户端广播的消息通道
    broadcast chan *models.Message

    // 注册请求通道
    register chan *Client

    // 注销请求通道
    unregister chan *Client

    // Redis客户端
    redisClient *redis.Client

    // 互斥锁
    mu sync.Mutex
}

// NewHub 创建一个新的Hub
func NewHub(redisClient *redis.Client) *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        users:       make(map[string]*models.User),
        broadcast:   make(chan *models.Message),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        redisClient: redisClient,
    }
}

// Run 启动Hub并处理通道消息
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)
        case client := <-h.unregister:
            h.unregisterClient(client)
        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}

// 注册一个新客户端
func (h *Hub) registerClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    // 创建用户数据
    user := &models.User{
        ID:         client.id,
        Nickname:   generateNickname(),
        JoinTime:   time.Now().UnixMilli(),
        LastActive: time.Now().UnixMilli(),
    }
    
    // 注册客户端和用户
    h.clients[client] = true
    h.users[client.id] = user
    
    // 保存用户数据到Redis
    h.redisClient.SaveUser(user)
    
    // 更新在线人数
    count := len(h.clients)
    h.redisClient.SetOnlineCount(count)
    
    // 发送欢迎消息
    welcomeMsg := &models.Message{
        ID: generateID(),
        Sender: models.User{
            ID:       "system",
            Nickname: "系统",
        },
        Content:   "欢迎来到聊天室",
        Timestamp: time.Now().UnixMilli(),
        Type:      "system",
    }
    client.send <- welcomeMsg
    
    // 发送用户信息和在线人数
    userInfo := &models.Message{
        ID:        generateID(),
        Content:   "user:info",
        Timestamp: time.Now().UnixMilli(),
        Type:      "user:info",
        Sender:    *user,
    }
    client.send <- userInfo
    
    // 发送聊天历史
    history, err := h.redisClient.GetChatHistory(50)
    if err == nil && len(history) > 0 {
        historyMsg := &models.Message{
            ID:        generateID(),
            Content:   "chat:history",
            Timestamp: time.Now().UnixMilli(),
            Type:      "chat:history",
            // 在前端处理时解析Data字段
            Data:      history,
        }
        client.send <- historyMsg
    }
    
    // 广播用户加入消息
    joinMsg := &models.Message{
        ID:        generateID(),
        Content:   user.Nickname + " 加入了聊天室",
        Timestamp: time.Now().UnixMilli(),
        Type:      "user:join",
        Sender:    *user,
    }
    h.broadcast <- joinMsg
    
    // 广播在线人数更新
    countMsg := &models.Message{
        ID:        generateID(),
        Content:   "online:count",
        Timestamp: time.Now().UnixMilli(),
        Type:      "online:count",
        Data:      count,
    }
    h.broadcast <- countMsg
}

// 注销一个客户端
func (h *Hub) unregisterClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if _, ok := h.clients[client]; ok {
        // 获取用户信息
        user, exists := h.users[client.id]
        
        // 删除客户端连接
        delete(h.clients, client)
        delete(h.users, client.id)
        close(client.send)
        
        // 从Redis中删除用户
        h.redisClient.RemoveUser(client.id)
        
        // 更新在线人数
        count := len(h.clients)
        h.redisClient.SetOnlineCount(count)
        
        if exists {
            // 广播用户离开消息
            leaveMsg := &models.Message{
                ID:        generateID(),
                Content:   user.Nickname + " 离开了聊天室",
                Timestamp: time.Now().UnixMilli(),
                Type:      "user:leave",
                Sender:    *user,
            }
            h.broadcast <- leaveMsg
        }
        
        // 广播在线人数更新
        countMsg := &models.Message{
            ID:        generateID(),
            Content:   "online:count",
            Timestamp: time.Now().UnixMilli(),
            Type:      "online:count",
            Data:      count,
        }
        h.broadcast <- countMsg
    }
}

// 广播消息给所有客户端
func (h *Hub) broadcastMessage(message *models.Message) {
    // 保存消息到Redis (非系统消息)
    if message.Type == "text" {
        h.redisClient.SaveChatMessage(message)
    }
    
    // 广播消息给所有客户端
    for client := range h.clients {
        select {
        case client.send <- message:
        default:
            // 如果客户端的发送缓冲区满了，关闭连接
            h.unregisterClient(client)
        }
    }
}
```

### 2.2 WebSocket客户端

每个WebSocket连接由一个Client对象表示：

```go
// internal/websocket/client.go

package websocket

import (
    "log"
    "time"
    
    "github.com/gorilla/websocket"
    "chat-server/internal/models"
    "chat-server/internal/utils"
)

const (
    // 写入超时
    writeWait = 10 * time.Second

    // 读取超时
    pongWait = 60 * time.Second

    // ping间隔
    pingPeriod = (pongWait * 9) / 10

    // 最大消息大小
    maxMessageSize = 1024 * 8
)

// Client 是WebSocket连接的中间人
type Client struct {
    // Hub实例
    hub *Hub

    // WebSocket连接
    conn *websocket.Conn

    // 缓冲的发送消息通道
    send chan *models.Message

    // 客户端唯一ID
    id string
}

// NewClient 创建一个新的客户端
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        hub:  hub,
        conn: conn,
        send: make(chan *models.Message, 256),
        id:   utils.GenerateID(),
    }
}

// 读取pump从WebSocket连接读取消息并将它们发送到hub
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    
    for {
        // 读取消息
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, 
                websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        
        // 解析消息并处理
        processedMessage := processMessage(c, message)
        if processedMessage != nil {
            c.hub.broadcast <- processedMessage
            
            // 更新用户最后活跃时间
            if user, ok := c.hub.users[c.id]; ok {
                user.LastActive = time.Now().UnixMilli()
                c.hub.redisClient.UpdateUserActivity(c.id, user.LastActive)
            }
        }
    }
}

// 写入pump将消息从hub发送到WebSocket连接
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Hub关闭了通道
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // 发送消息
            err := c.conn.WriteJSON(message)
            if err != nil {
                log.Printf("error sending message: %v", err)
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

### 2.3 Redis客户端

Redis客户端管理与Redis的通信：

```go
// internal/redis/client.go

package redis

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    
    "github.com/go-redis/redis/v8"
    "chat-server/internal/models"
)

// 常量
const (
    chatHistoryKey = "chat:history"
    usersKey       = "chat:users"
    onlineCountKey = "chat:onlineCount"
    maxHistorySize = 100
)

// Client 封装了Redis操作
type Client struct {
    rdb *redis.Client
    ctx context.Context
}

// NewClient 创建新的Redis客户端
func NewClient(addr string, password string, db int) (*Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    ctx := context.Background()
    
    // 测试连接
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }
    
    return &Client{rdb: rdb, ctx: ctx}, nil
}

// SaveChatMessage 保存聊天消息到Redis
func (c *Client) SaveChatMessage(message *models.Message) error {
    // 将消息序列化为JSON
    jsonData, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    // 使用管道执行多个命令
    pipe := c.rdb.Pipeline()
    
    // 添加消息到列表开头
    pipe.LPush(c.ctx, chatHistoryKey, jsonData)
    
    // 限制列表大小
    pipe.LTrim(c.ctx, chatHistoryKey, 0, maxHistorySize-1)
    
    // 执行管道
    _, err = pipe.Exec(c.ctx)
    return err
}

// GetChatHistory 获取聊天历史记录
func (c *Client) GetChatHistory(limit int) ([]*models.Message, error) {
    if limit <= 0 {
        limit = maxHistorySize
    }
    
    // 从Redis获取历史记录
    jsonMessages, err := c.rdb.LRange(c.ctx, chatHistoryKey, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, err
    }
    
    // 解析消息
    messages := make([]*models.Message, 0, len(jsonMessages))
    for _, jsonMsg := range jsonMessages {
        var msg models.Message
        if err := json.Unmarshal([]byte(jsonMsg), &msg); err != nil {
            continue
        }
        messages = append(messages, &msg)
    }
    
    // 反转消息顺序
    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }
    
    return messages, nil
}

// SaveUser 保存用户信息到Redis
func (c *Client) SaveUser(user *models.User) error {
    jsonData, err := json.Marshal(user)
    if err != nil {
        return err
    }
    
    return c.rdb.HSet(c.ctx, usersKey, user.ID, jsonData).Err()
}

// RemoveUser 从Redis中删除用户
func (c *Client) RemoveUser(userID string) error {
    return c.rdb.HDel(c.ctx, usersKey, userID).Err()
}

// GetAllUsers 获取所有在线用户
func (c *Client) GetAllUsers() ([]*models.User, error) {
    jsonUsers, err := c.rdb.HGetAll(c.ctx, usersKey).Result()
    if err != nil {
        return nil, err
    }
    
    users := make([]*models.User, 0, len(jsonUsers))
    for _, jsonUser := range jsonUsers {
        var user models.User
        if err := json.Unmarshal([]byte(jsonUser), &user); err != nil {
            continue
        }
        users = append(users, &user)
    }
    
    return users, nil
}

// UpdateUserActivity 更新用户最后活跃时间
func (c *Client) UpdateUserActivity(userID string, timestamp int64) error {
    // 获取用户数据
    jsonUser, err := c.rdb.HGet(c.ctx, usersKey, userID).Result()
    if err != nil {
        if err == redis.Nil {
            return errors.New("用户不存在")
        }
        return err
    }
    
    var user models.User
    if err := json.Unmarshal([]byte(jsonUser), &user); err != nil {
        return err
    }
    
    // 更新活跃时间
    user.LastActive = timestamp
    
    // 保存回Redis
    return c.SaveUser(&user)
}

// SetOnlineCount 设置在线人数
func (c *Client) SetOnlineCount(count int) error {
    return c.rdb.Set(c.ctx, onlineCountKey, count, 0).Err()
}

// GetOnlineCount 获取在线人数
func (c *Client) GetOnlineCount() (int, error) {
    val, err := c.rdb.Get(c.ctx, onlineCountKey).Int()
    if err != nil {
        if err == redis.Nil {
            return 0, nil
        }
        return 0, err
    }
    return val, nil
}
```

## 3. HTTP API实现

REST API使用Gin框架实现：

```go
// internal/api/handlers.go

package api

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "chat-server/internal/redis"
)

// Handler 包含所有HTTP处理器
type Handler struct {
    redisClient *redis.Client
}

// NewHandler 创建一个新的API处理器
func NewHandler(redisClient *redis.Client) *Handler {
    return &Handler{
        redisClient: redisClient,
    }
}

// GetChatHistory 处理获取聊天历史的请求
func (h *Handler) GetChatHistory(c *gin.Context) {
    // 解析分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
    
    // 计算起始索引
    start := page * pageSize
    limit := start + pageSize
    
    // 获取历史记录
    messages, err := h.redisClient.GetChatHistory(limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取聊天历史失败"})
        return
    }
    
    // 如果结果集小于起始索引，返回空数组
    if len(messages) <= start {
        c.JSON(http.StatusOK, gin.H{"messages": []interface{}{}})
        return
    }
    
    // 返回指定页的消息
    end := len(messages)
    if end > limit {
        end = limit
    }
    
    c.JSON(http.StatusOK, gin.H{"messages": messages[start:end]})
}

// GetUsers 处理获取用户列表的请求
func (h *Handler) GetUsers(c *gin.Context) {
    users, err := h.redisClient.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetOnlineCount 处理获取在线人数的请求
func (h *Handler) GetOnlineCount(c *gin.Context) {
    count, err := h.redisClient.GetOnlineCount()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取在线人数失败"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"count": count})
}
```

```go
// internal/api/routes.go

package api

import (
    "github.com/gin-gonic/gin"
    "chat-server/internal/websocket"
    "chat-server/internal/redis"
)

// SetupRoutes 配置所有路由
func SetupRoutes(router *gin.Engine, hub *websocket.Hub, redisClient *redis.Client) {
    // 创建处理器
    handler := NewHandler(redisClient)
    
    // API路由组
    api := router.Group("/api")
    {
        // 聊天相关API
        chat := api.Group("/chat")
        {
            chat.GET("/history", handler.GetChatHistory)
            chat.GET("/users", handler.GetUsers)
            chat.GET("/online-count", handler.GetOnlineCount)
        }
    }
    
    // WebSocket路由
    router.GET("/ws", func(c *gin.Context) {
        websocket.ServeWs(hub, c.Writer, c.Request)
    })
}
```

## 4. 主程序入口

```go
// cmd/server/main.go

package main

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    
    "chat-server/internal/api"
    "chat-server/internal/config"
    "chat-server/internal/redis"
    "chat-server/internal/websocket"
)

func main() {
    // 解析命令行参数
    configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
    flag.Parse()
    
    // 加载配置
    cfg, err := config.LoadConfig(*configPath)
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 连接Redis
    redisClient, err := redis.NewClient(
        cfg.Redis.Addr,
        cfg.Redis.Password,
        cfg.Redis.DB,
    )
    if err != nil {
        log.Fatalf("连接Redis失败: %v", err)
    }
    
    // 创建WebSocket hub
    hub := websocket.NewHub(redisClient)
    go hub.Run()
    
    // 创建Gin实例
    router := gin.Default()
    
    // 配置CORS
    router.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.CORS.AllowOrigins,
        AllowMethods:     cfg.CORS.AllowMethods,
        AllowHeaders:     cfg.CORS.AllowHeaders,
        AllowCredentials: true,
    }))
    
    // 设置路由
    api.SetupRoutes(router, hub, redisClient)
    
    // 启动HTTP服务器
    go func() {
        if err := router.Run(cfg.Server.Addr); err != nil {
            log.Fatalf("启动服务器失败: %v", err)
        }
    }()
    
    log.Printf("服务器已启动，监听地址: %s", cfg.Server.Addr)
    
    // 等待中断信号优雅地关闭服务器
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("正在关闭服务器...")
}
```

## 5. 配置结构

```go
// internal/config/config.go

package config

import (
    "os"
    
    "gopkg.in/yaml.v2"
)

// Config 应用程序配置
type Config struct {
    Server ServerConfig `yaml:"server"`
    Redis  RedisConfig  `yaml:"redis"`
    CORS   CORSConfig   `yaml:"cors"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
    Addr string `yaml:"addr"`
}

// RedisConfig Redis配置
type RedisConfig struct {
    Addr     string `yaml:"addr"`
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
}

// CORSConfig CORS配置
type CORSConfig struct {
    AllowOrigins []string `yaml:"allow_origins"`
    AllowMethods []string `yaml:"allow_methods"`
    AllowHeaders []string `yaml:"allow_headers"`
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

## 6. 示例配置文件

```yaml
# configs/config.yaml

server:
  addr: ":8080"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

cors:
  allow_origins:
    - "http://localhost:5173"
  allow_methods:
    - "GET"
    - "POST"
  allow_headers:
    - "Content-Type"
    - "Authorization"
```

## 7. 实现细节和注意事项

### 7.1 消息处理

WebSocket消息处理逻辑：

```go
// internal/websocket/message.go
func processMessage(client *Client, data []byte) *models.Message {
    // 解析消息类型和内容
    var rawMsg map[string]interface{}
    if err := json.Unmarshal(data, &rawMsg); err != nil {
        return nil
    }
    
    // 根据消息类型处理
    msgType, ok := rawMsg["type"].(string)
    if !ok {
        return nil
    }
    
    // 获取用户
    user, ok := client.hub.users[client.id]
    if !ok {
        return nil
    }
    
    switch msgType {
    case "chat:message":
        // 处理聊天消息
        content, ok := rawMsg["content"].(string)
        if !ok || content == "" {
            return nil
        }
        
        // 过滤和清理内容
        content = utils.SanitizeInput(content)
        
        // 检查消息长度
        if len(content) > 500 {
            content = content[:500]
        }
        
        // 创建消息对象
        return &models.Message{
            ID:        utils.GenerateID(),
            Sender:    *user,
            Content:   content,
            Timestamp: time.Now().UnixMilli(),
            Type:      "text",
        }
        
    case "user:nickname":
        // 处理昵称更新
        nickname, ok := rawMsg["content"].(string)
        if !ok || nickname == "" {
            return nil
        }
        
        // 清理和检查昵称
        nickname = utils.SanitizeInput(nickname)
        if len(nickname) > 20 {
            nickname = nickname[:20]
        }
        
        // 更新用户昵称
        oldNickname := user.Nickname
        user.Nickname = nickname
        
        // 保存到Redis
        client.hub.redisClient.SaveUser(user)
        
        // 创建昵称更新消息
        return &models.Message{
            ID:        utils.GenerateID(),
            Content:   oldNickname + " 改名为 " + nickname,
            Timestamp: time.Now().UnixMilli(),
            Type:      "user:nickname",
            Sender:    *user,
        }
    }
    
    return nil
}
```

### 7.2 安全工具

```go
// internal/utils/security.go

package utils

import (
    "html"
    "regexp"
    "strings"
)

var scriptPattern = regexp.MustCompile(`(?i)<script[\s\S]*?>[\s\S]*?</script>`)
var tagPattern = regexp.MustCompile(`(?i)<[^>]*>`)

// SanitizeInput 清理输入以防止XSS
func SanitizeInput(input string) string {
    // 去除script标签及其内容
    noScripts := scriptPattern.ReplaceAllString(input, "")
    
    // 去除所有HTML标签
    noTags := tagPattern.ReplaceAllString(noScripts, "")
    
    // 转义HTML实体
    escaped := html.EscapeString(noTags)
    
    // 去除多余空白
    return strings.TrimSpace(escaped)
}
```

### 7.3 ID生成

```go
// internal/utils/id.go

package utils

import (
    "math/rand"
    "strconv"
    "time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateID 生成一个唯一ID
func GenerateID() string {
    timestamp := strconv.FormatInt(time.Now().UnixNano(), 36)
    
    // 添加5位随机字符
    b := make([]byte, 5)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    
    return timestamp + string(b)
}

// 生成临时昵称
func generateNickname() string {
    return "游客" + strconv.Itoa(seededRand.Intn(10000))
}
```

## 8. 构建和运行

### 8.1 依赖管理

使用Go模块管理依赖：

```bash
# 初始化模块
go mod init chat-server

# 安装依赖
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/gorilla/websocket
go get github.com/go-redis/redis/v8
go get gopkg.in/yaml.v2
```

### 8.2 构建和运行

```bash
# 构建
go build -o chat-server ./cmd/server

# 运行
./chat-server
```

### 8.3 Docker支持

创建Dockerfile：

```dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chat-server .
COPY configs/ configs/

EXPOSE 8080

CMD ["./chat-server"]
``` 