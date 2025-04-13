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
	Server    ServerConfig    `yaml:"server"`
	Redis     RedisConfig     `yaml:"redis"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	Auth      AuthConfig      `yaml:"auth"`
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
	Path             string `yaml:"path"`
	MessageSizeLimit int64  `yaml:"message_size_limit"`
	MessageBufferSize int    `yaml:"message_buffer_size"`
}

// AuthConfig 包含认证相关配置
type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry string `yaml:"token_expiry"`
}

// 全局配置实例
var AppConfig Config

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

	return nil
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