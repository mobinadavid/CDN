package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisService struct {
	RedisClient *redis.Client
}

func NewRedisService(redisClient *redis.Client) *RedisService {
	return &RedisService{RedisClient: redisClient}
}

func (redisService *RedisService) CheckRateLimit(ip string) (bool, error) {
	// Get the current count of objects put by the IP
	count, err := redisService.RedisClient.Get(context.Background(), ip).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// If count doesn't exist, initialize it to 0
	if err == redis.Nil {
		count = 0
	}

	// If count exceeds 10, reject the request
	if count >= 10 {
		return false, nil
	}

	// Increment the count
	err = redisService.RedisClient.Incr(context.Background(), ip).Err()
	if err != nil {
		return false, err
	}

	// Set the expiration time for the key to one hour if it doesn't exist
	redisService.RedisClient.Expire(context.Background(), ip, time.Hour)

	return true, nil
}
