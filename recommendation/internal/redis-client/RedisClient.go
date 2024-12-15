package redis_client

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

// redisClient 초기화 메서드
func NewRedisClient(addr string, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
