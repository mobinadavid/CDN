package drivers

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
	Client *redis.Client
)

// Connect establishes a connection to Redis
func (r *Redis) Connect() (err error) {
	//connStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s")

	// Create Redis client
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       0, //todo: should reconsider.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	})

	// Ping Redis server to ensure connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = Client.Ping(ctx).Result(); err != nil {
		log.Fatalln(err)
	}

	return
}

// Close closes the connection to Redis
func (r *Redis) Close() (err error) {
	if err = Client.Close(); err != nil {
		log.Fatalln(err)
	}
	return
}

func (r *Redis) GetClient() *redis.Client {
	return Client
}

func (r *Redis) CheckRateLimit(ip string) (bool, error) {
	// Get the current count of objects put by the IP
	count, err := r.GetClient().Get(context.Background(), ip).Int()
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
	err = r.GetClient().Incr(context.Background(), ip).Err()
	if err != nil {
		return false, err
	}

	// Set the expiration time for the key to one hour if it doesn't exist
	r.GetClient().Expire(context.Background(), ip, time.Hour)

	return true, nil
}

func Init() error {
	configs := config.GetInstance()
	redis := Redis{
		Host:     configs.Get("REDIS_HOST"),
		Port:     configs.Get("REDIS_PORT"),
		Password: configs.Get("REDIS_PASSWORD"),
		Database: 0,
	}
	err := redis.Connect()
	return err
}
