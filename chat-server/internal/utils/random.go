package utils

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	// 确保随机数生成器已初始化
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rnd.Intn(len(letterBytes))]
	}
	return string(b)
} 