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

func (objectService *ObjectService) PutObject(ctx context.Context, bucket string, files []*multipart.FileHeader, folder string, tags ...map[string]string) ([]map[string]string, error) {

	var uploadInfoList []map[string]string

	for _, file := range files {

		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		fileName := utils.GenerateUUIDFileName(file.Filename)
		uuidFileName := folder + "/" + fileName

		options := minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		}

		if len(tags) != 0 && tags[0] != nil {
			options.UserTags = tags[0]
		}

		_, err = objectService.MinioClient.PutObject(ctx, bucket, uuidFileName, src, file.Size, options)
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

func (objectService *ObjectService) GetObject(ctx context.Context, bucket, fileName string, options minio.GetObjectOptions) (*minio.Object, error) {

	return objectService.MinioClient.GetObject(ctx, bucket, fileName, options)
}

func (objectService *ObjectService) RemoveObjects(ctx context.Context, bucketName, objectName string, options minio.RemoveObjectOptions) error {

	err := objectService.MinioClient.RemoveObject(ctx, bucketName, objectName, options)
	return err

}

func (objectService *ObjectService) GetObjectsByTags(ctx context.Context, bucketName string, tags map[string]string) ([]map[string]string, error) {

	var uploadInfoList []map[string]string
	existingTag := make(map[string]string)

	// List objects in the bucket
	objectCh := objectService.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		tag, err := objectService.MinioClient.GetObjectTagging(ctx, bucketName, object.Key, minio.GetObjectTaggingOptions{})
		if err != nil {
			return nil, err
		}

		if tag.String() != "" {
			tagPairs := strings.Split(tag.String(), "&")
			for _, tagPair := range tagPairs {
				pair := strings.Split(tagPair, "=")
				if len(pair) == 2 {
					existingTag[pair[0]] = pair[1]
				}
			}
		}

		if tagsMatch(existingTag, tags) {
			uploadInfoList = append(uploadInfoList, map[string]string{
				"url": fmt.Sprintf("%objectService://%objectService/%objectService/%objectService/%objectService",
					ctx.Value("Scheme"),
					ctx.Value("Host"),
					"app/api/v1/storage",
					bucketName,
					object.Key,
				),
			})
		}
	}

	return uploadInfoList, nil
}

func tagsMatch(objectTags map[string]string, queryTags map[string]string) bool {

	for key, value := range queryTags {
		if objectTags[key] != value {
			return false
		}
	}
	return true
}
