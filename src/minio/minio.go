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
	clientCDN      *minio.Client
	clientInternal *minio.Client
}

// Init initializes the MinIO clients.
func Init() error {
	return GetInstance().Connect()
}

// Connect connects to MinIO.
func (m *Client) Connect() error {
	configs = config.GetInstance()

	// Create the internal MinIO client
	internalClient, err := minio.New(fmt.Sprintf("%s:%s", configs.Get("MINIO_HOST"), configs.Get("MINIO_PORT_INTERNAL")), &minio.Options{
		Creds:  credentials.NewStaticV4(configs.Get("MINIO_ACCESS_KEY"), configs.Get("MINIO_SECRET_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		return err
	}
	m.clientInternal = internalClient

	// Create the CDN MinIO client
	cdnClient, err := minio.New(fmt.Sprintf("%s", configs.Get("MINIO_CDN_HOST")), &minio.Options{
		Creds:  credentials.NewStaticV4(configs.Get("MINIO_ACCESS_KEY"), configs.Get("MINIO_SECRET_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		return err
	}
	m.clientCDN = cdnClient

	return nil
}

// GetCDNClient returns the CDN MinIO client.
func (m *Client) GetCDNClient() *minio.Client {
	return m.clientCDN
}

// GetInternalClient returns the internal MinIO client.
func (m *Client) GetInternalClient() *minio.Client {
	return m.clientInternal
}

// GetInstance returns the singleton instance of Minio.
func GetInstance() *Client {
	if instance == nil {
		instance = &Client{}
	}
	return instance
}
