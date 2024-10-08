package minio

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type BucketService struct {
	MinioClient *minio.Client
}

func NewBucketService(minioClient *minio.Client) *BucketService {
	return &BucketService{MinioClient: minioClient}
}

func (bucketService *BucketService) MakeBucket(ctx context.Context, name string, options minio.MakeBucketOptions) error {
	return bucketService.MinioClient.MakeBucket(ctx, name, options)
}

func (bucketService *BucketService) ListBucket(ctx context.Context) ([]minio.BucketInfo, error) {
	return bucketService.MinioClient.ListBuckets(ctx)
}

func (bucketService *BucketService) RemoveBucket(ctx context.Context, name string) error {
	return bucketService.MinioClient.RemoveBucket(ctx, name)
}

func (bucketService *BucketService) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return bucketService.MinioClient.BucketExists(ctx, bucket)
}

func (bucketService *BucketService) ListObjects(ctx context.Context, bucketName string, options minio.ListObjectsOptions) ([]minio.ObjectInfo, error) {
	exists, err := bucketService.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("bucket does not exist")
	}

	var objects []minio.ObjectInfo

	objectCh := bucketService.MinioClient.ListObjects(ctx, bucketName, options)
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error occurred while listing objects: %w", object.Err)

		}
		objects = append(objects, object)
	}

	return objects, nil
}
