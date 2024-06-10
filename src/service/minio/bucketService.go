package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type BucketService struct {
	MinioClient *minio.Client
}

func NewBucketService(minioClient *minio.Client) *BucketService {
	return &BucketService{MinioClient: minioClient}
}

func (bucketSerevice *BucketService) MakeBucket(ctx context.Context, name string, options minio.MakeBucketOptions) error {
	return bucketSerevice.MinioClient.MakeBucket(ctx, name, options)
}
func (storageService *BucketService) RemoveBucket(ctx context.Context, name string) error {
	return storageService.MinioClient.RemoveBucket(ctx, name)
}

func (storageService *BucketService) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return storageService.MinioClient.BucketExists(ctx, bucket)
}
func (storageService *BucketService) ListObjects(ctx context.Context, bucketName string, options minio.ListObjectsOptions) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo
	objectCh := storageService.MinioClient.ListObjects(ctx, bucketName, options)
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error occurred while listing objects: %w", object.Err)

		}
		objects = append(objects, object)
	}

	return objects, nil
}
