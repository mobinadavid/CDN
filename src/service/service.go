package service

import (
	"cdn/src/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/context"
	"mime/multipart"
	"strconv"
	"strings"
)

type StorageService struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewStorageService(minioClient *minio.Client, bucketName string) *StorageService {
	return &StorageService{
		MinioClient: minioClient,
		BucketName:  bucketName,
	}
}

func (s *StorageService) UploadFile(c *gin.Context, file *multipart.FileHeader) (map[string]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	uuidFileName := utils.GenerateUUIDFileName(file.Filename)
	_, err = s.MinioClient.PutObject(c, s.BucketName, uuidFileName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"original_file_name": strings.ToLower(file.Filename),
		"size":               strconv.FormatInt(file.Size, 10),
		"file_name":          uuidFileName,
		"url": fmt.Sprintf("%s://%s/%s/%s",
			c.GetHeader("Scheme"),
			c.Request.Host,
			"app/api/v1/storage",
			uuidFileName,
		),
	}, nil
}
func (s *StorageService) GetObject(fileName string) (*minio.Object, error) {
	return s.MinioClient.GetObject(context.Background(), s.BucketName, fileName, minio.GetObjectOptions{})
}
