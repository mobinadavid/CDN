package minio

import (
	"cdn/src/pkg/utils"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"strconv"
	"strings"
)

type ObjectService struct {
	MinioClient *minio.Client
}

func NewObjectService(minioClient *minio.Client) *ObjectService {
	return &ObjectService{MinioClient: minioClient}
}
func (objectService *ObjectService) PutObject(ctx context.Context, bucket string, files []*multipart.FileHeader, folder string) ([]map[string]string, error) {
	var uploadInfoList []map[string]string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		fileName := utils.GenerateUUIDFileName(file.Filename)

		uuidFileName := folder + "/" + fileName

		_, err = objectService.MinioClient.PutObject(ctx, bucket, uuidFileName, src, file.Size, minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
		if err != nil {
			return nil, err
		}

		uploadInfoList = append(uploadInfoList, map[string]string{
			"original_file_name": strings.ToLower(file.Filename),
			"size":               strconv.FormatInt(file.Size, 10),
			"file_name":          fileName,
			"folder":             folder,
			"url": fmt.Sprintf("%objectService://%objectService/%objectService/%objectService/%objectService",
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
func (storageService *ObjectService) GetObject(ctx context.Context, bucket, fileName string, options minio.GetObjectOptions) (*minio.Object, error) {
	return storageService.MinioClient.GetObject(ctx, bucket, fileName, options)
}

func (storageService *ObjectService) RemoveObjects(ctx context.Context, bucketName, objectName string, options minio.RemoveObjectOptions) error {
	err := storageService.MinioClient.RemoveObject(ctx, bucketName, objectName, options)
	return err

}
