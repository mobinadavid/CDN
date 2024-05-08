package minio

import (
	"cdn/src/config"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	configs  *config.Config
	instance *Client
)

// Client represents a client for interacting with Minio.
type Client struct {
	client *minio.Client
}

// Init creates a new Minio client.
func Init() error {
	return GetInstance().Connect()
}

// Connect Connects to Minio.
func (m *Client) Connect() (err error) {
	configs = config.GetInstance()

	m.client, err = minio.New(fmt.Sprintf("%s:9000", configs.Get("MINIO_HOST")), &minio.Options{
		Creds:  credentials.NewStaticV4(configs.Get("MINIO_ACCESS_KEY"), configs.Get("MINIO_SECRET_KEY"), ""),
		Secure: true,
	})

	if err != nil {
		return err
	}

	return nil
}

// GetMinio returns the singleton instance of Minio.
func (m *Client) GetMinio() *minio.Client {
	return m.client
}

// GetInstance returns the singleton instance of Minio.
func GetInstance() *Client {
	if instance == nil {
		instance = &Client{}
	}
	return instance
}
