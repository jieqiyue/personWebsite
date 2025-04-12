package utils

import (
	"github.com/google/uuid"
)

// GenerateUUID 生成一个UUID字符串
func GenerateUUID() string {
	return uuid.New().String()
} 