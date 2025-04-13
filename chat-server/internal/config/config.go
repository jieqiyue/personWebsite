package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config 保存服务器的所有配置信息
type Config struct {
	Server    ServerConfig     `yaml:"server"`
	Redis     RedisConfig      `yaml:"redis"`
	WebSocket WebSocketConfig  `yaml:"websocket"`
	Auth      AuthConfig       `yaml:"auth"`
	ChatRooms []ChatRoomConfig `yaml:"chat_rooms"` // 聊天室配置
}

// ServerConfig 包含HTTP服务器配置
type ServerConfig struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

// RedisConfig 包含Redis服务器配置
type RedisConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	Prefix      string `yaml:"prefix"`
	MaxMessages int64  `yaml:"max_messages"`
}

// WebSocketConfig 包含WebSocket配置
type WebSocketConfig struct {
	Path              string `yaml:"path"`
	MessageSizeLimit  int64  `yaml:"message_size_limit"`
	MessageBufferSize int    `yaml:"message_buffer_size"`
}

// AuthConfig 包含认证相关配置
type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry string `yaml:"token_expiry"`
}

// ChatRoomConfig 包含聊天室配置
type ChatRoomConfig struct {
	ID          string `yaml:"id"`          // 房间ID
	Name        string `yaml:"name"`        // 房间名称
	Description string `yaml:"description"` // 房间描述
	IsDefault   bool   `yaml:"is_default"`  // 是否为默认房间
}

// 全局配置实例
var AppConfig Config

// 聊天室ID集合，用于快速验证
var ValidRoomIDs map[string]bool

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) error {
	configFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 初始化有效房间ID映射
	initValidRoomIDs()

	return nil
}

// initValidRoomIDs 初始化有效房间ID映射，用于快速查找
func initValidRoomIDs() {
	ValidRoomIDs = make(map[string]bool)
	for _, room := range AppConfig.ChatRooms {
		ValidRoomIDs[room.ID] = true
	}
}

// IsValidRoomID 判断给定的房间ID是否有效
func IsValidRoomID(roomID string) bool {
	return ValidRoomIDs[roomID]
}

// GetDefaultRoomID 获取默认房间ID
func GetDefaultRoomID() string {
	for _, room := range AppConfig.ChatRooms {
		if room.IsDefault {
			return room.ID
		}
	}
	// 如果没有设置默认房间，返回第一个房间ID，若无房间则返回"general"
	//if len(AppConfig.ChatRooms) > 0 {
	//	return AppConfig.ChatRooms[0].ID
	//}
	return "general"
}

// GetDefaultConfigPath 返回默认配置文件路径
func GetDefaultConfigPath() string {
	// 尝试从当前目录或父目录找到configs/config.yaml
	currentDir, _ := os.Getwd()
	for i := 0; i < 3; i++ { // 向上查找最多3层
		path := filepath.Join(currentDir, "configs", "config.yaml")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		currentDir = filepath.Dir(currentDir)
	}

	// 没找到返回默认路径
	return "configs/config.yaml"
}
