package service

import (
	"cdn/src/pkg/utils"
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
)

type StorageService struct {
	MinioClient *minio.Client
}

func NewStorageService(minioClient *minio.Client) *StorageService {
	return &StorageService{MinioClient: minioClient}
}

func (storageService *StorageService) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return storageService.MinioClient.BucketExists(ctx, bucket)
}

func (storageService *StorageService) UploadFiles(ctx context.Context, bucket string, files []*multipart.FileHeader) ([]map[string]string, error) {
	var uploadInfoList []map[string]string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		uuidFileName := utils.GenerateUUIDFileName(file.Filename)

		_, err = storageService.MinioClient.PutObject(ctx, bucket, uuidFileName, src, file.Size, minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
		if err != nil {
			return nil, err
		}

		uploadInfoList = append(uploadInfoList, map[string]string{
			"original_file_name": strings.ToLower(file.Filename),
			"size":               strconv.FormatInt(file.Size, 10),
			"file_name":          uuidFileName,
			"url": fmt.Sprintf("%storageService://%storageService/%storageService/%storageService/%storageService",
				ctx.Value("Scheme"),
				ctx.Value("Host"),
				"app/api/v1/storage",
				bucket,
				uuidFileName,
			),
		})
	}

	return uploadInfoList, nil
}
func (storageService *StorageService) GetObject(ctx context.Context, bucket, fileName string, options minio.GetObjectOptions) (*minio.Object, error) {
	return storageService.MinioClient.GetObject(ctx, bucket, fileName, options)
}
