package redis

import (
	"cdn/src/config"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Redis struct {
	Host     string
	Port     string
	Password string
	Database int
}

var (
	configs  *config.Config
	instance *Client
)

type Client struct {
	client *redis.Client
}

// Connect establishes a connection to Redis
func (r *Client) Connect() (err error) {
	configs = config.GetInstance()
	// Create Redis client
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.Get("REDIS_HOST"), configs.Get("REDIS_PORT")),
		Password: configs.Get("REDIS_PASSWORD"),
		DB:       0, //todo: should reconsider.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	})

	// Ping Redis server to ensure connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = r.client.Ping(ctx).Result(); err != nil {
		log.Fatalln(err)
	}

	return
}

// Close closes the connection to Redis
func (r *Client) Close() (err error) {
	if err = r.Close(); err != nil {
		log.Fatalln(err)
	}
	return
}

func (r *Client) GetClient() *redis.Client {
	return r.client
}
func GetInstance() *Client {
	if instance == nil {
		instance = &Client{}
	}
	return instance
}

//	func (r *Redis) CheckRateLimit(ip string) (bool, error) {
//		// Get the current count of objects put by the IP
//		count, err := r.GetClient().Get(context.Background(), ip).Int()
//		if err != nil && err != redis.Nil {
//			return false, err
//		}
//
//		// If count doesn't exist, initialize it to 0
//		if err == redis.Nil {
//			count = 0
//		}
//
//		// If count exceeds 10, reject the request
//		if count >= 10 {
//			return false, nil
//		}
//
//		// Increment the count
//		err = r.GetClient().Incr(context.Background(), ip).Err()
//		if err != nil {
//			return false, err
//		}
//
//		// Set the expiration time for the key to one hour if it doesn't exist
//		r.GetClient().Expire(context.Background(), ip, time.Hour)
//
//		return true, nil
//	}
func Init() error {
	return GetInstance().Connect()
}
