package redis

import (
	"cdn/src/config"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisService struct {
	RedisClient *redis.Client
}

func NewRedisService(redisClient *redis.Client) *RedisService {

	return &RedisService{RedisClient: redisClient}
}

func (redisService *RedisService) GenerateCompositeKey(ip, userAgent string) string {
	compositeKey := fmt.Sprintf("%s:%s", ip, userAgent)
	hasher := sha1.New()
	hasher.Write([]byte(compositeKey))

	return fmt.Sprintf("%x", hasher.Sum(nil))

}

func (redisService *RedisService) CheckAndIncrementRateLimit(ip, userAgent string) (bool, error) {
	configs := config.GetInstance()
	ctx := context.Background()
	compositeKey := redisService.GenerateCompositeKey(ip, userAgent)

	rateLimitStr := configs.Get("RATE_LIMIT")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		return false, err // Handle conversion error
	}
	// Get the current count of objects put by the composite key
	count, err := redisService.RedisClient.Get(ctx, compositeKey).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// If count doesn't exist, initialize it to 0
	if err == redis.Nil {
		count = 0
	}
	fmt.Println("")
	// If count exceeds the limit, reject the request
	if count >= rateLimit {
		return false, nil
	}

	// Increment the count
	err = redisService.RedisClient.Incr(ctx, compositeKey).Err()
	if err != nil {
		return false, err
	}

	// Set the expiration time for the key to one hour if it doesn't exist
	if count == 0 {
		err = redisService.RedisClient.Expire(ctx, compositeKey, time.Hour).Err()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
