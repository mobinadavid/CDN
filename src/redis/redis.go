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

func Init() error {
	return GetInstance().Connect()
}
