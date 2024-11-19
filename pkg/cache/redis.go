package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ctx = context.Background()

type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(dsn string) *RedisTokenStore {
	client := redis.NewClient(&redis.Options{
		Addr: dsn,
		DB:   0, // 使用的数据库编号
	})

	// 测试连接，确保 Redis 可用
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("client.Ping in NewRedisTokenStore OK")

	return &RedisTokenStore{client: client}
}

// 设置 refresh token
func (rts *RedisTokenStore) Set(token, userID string, duration time.Duration) error {
	return rts.client.Set(ctx, token, userID, duration).Err()
}

func (rts *RedisTokenStore) Get(token string) (string, bool) {
	userID, err := rts.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return "", false // Token 不存在
	} else if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return "", false
	}
	return userID, true
}

func (rts *RedisTokenStore) Demo() {
	// 设置 token
	err := rts.Set("token123", "user123", 2*time.Minute)
	if err != nil {
		fmt.Printf("Error setting token: %v\n", err)
		return
	}

	// 获取 token
	userID, exists := rts.Get("token123")
	if exists {
		fmt.Printf("Token belongs to user: %s\n", userID)
	} else {
		fmt.Println("Token not found or expired")
	}
}
