package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"chat-server/internal/config"
	"chat-server/internal/models"
)

// Client Redis客户端实例
var Client *redis.Client

// Setup 初始化Redis连接
func Setup() error {
	cfg := config.AppConfig.Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// GetKey 获取带前缀的完整键名
func GetKey(key string) string {
	return config.AppConfig.Redis.Prefix + key
}

// SaveMessage 保存消息到Redis
func SaveMessage(ctx context.Context, roomID string, message []byte) error {
	key := GetKey("messages:" + roomID)

	// 添加消息到列表末尾
	err := Client.RPush(ctx, key, message).Err()
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// 使用LTRIM限制消息列表长度，只保留最新的消息
	// 从配置中获取最大消息数量
	maxMessages := config.AppConfig.Redis.MaxMessages
	if maxMessages <= 0 {
		maxMessages = 1000 // 默认值，以防配置文件中没有设置
	}

	// 使用负数作为起始索引，表示从列表末尾向前数
	err = Client.LTrim(ctx, key, -maxMessages, -1).Err()
	if err != nil {
		return fmt.Errorf("failed to trim message list: %w", err)
	}

	// 移除过期时间设置，让消息永久保存（只受数量限制）

	return nil
}

func GetDefaultMessageLen(ctx context.Context, roomID string) int {
	key := GetKey("messages:" + roomID)
	listLen, err := Client.LLen(ctx, key).Result()
	if err != nil {
		return 0
	}

	return int(listLen)
}

// GetMessages 获取特定房间的消息历史
func GetMessages(ctx context.Context, roomID string, limit, offset int64) ([]models.Message, error) {
	if GetDefaultMessageLen(ctx, roomID) == 0 {
		return nil, nil
	}

	key := GetKey("messages:" + roomID)

	messageStrs, err := Client.LRange(ctx, key, limit, offset).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	
	messages := make([]models.Message, 0, len(messageStrs))
	for _, msgStr := range messageStrs {
		var msg models.Message
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			// 跳过解析错误的消息
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SaveUser 保存用户信息到Redis
func SaveUser(ctx context.Context, userID string, userData []byte) error {
	key := GetKey("users:" + userID)

	// 设置30分钟过期，临时用户信息
	err := Client.Set(ctx, key, userData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// GetUser 从Redis获取用户信息
func GetUser(ctx context.Context, userID string) (*models.User, error) {
	key := GetKey("users:" + userID)

	userData, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user models.User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &user, nil
}

// AddUserToRoom 将用户添加到房间
func AddUserToRoom(ctx context.Context, roomID, userID string) error {
	key := GetKey("room_members:" + roomID)

	err := Client.SAdd(ctx, key, userID).Err()
	if err != nil {
		return fmt.Errorf("failed to add user to room: %w", err)
	}

	return nil
}

// RemoveUserFromRoom 将用户从房间移除
func RemoveUserFromRoom(ctx context.Context, roomID, userID string) error {
	key := GetKey("room_members:" + roomID)

	err := Client.SRem(ctx, key, userID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove user from room: %w", err)
	}

	return nil
}

// GetRoomMembers 获取房间中的所有用户
func GetRoomMembers(ctx context.Context, roomID string) ([]string, error) {
	key := GetKey("room_members:" + roomID)

	members, err := Client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get room members: %w", err)
	}

	return members, nil
}
