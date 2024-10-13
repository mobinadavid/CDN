package minio

import (
	"cdn/src/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"strconv"
	"strings"
)

type ObjectService struct {
	MinioClient   *minio.Client
	BucketService *BucketService
}

func NewObjectService(minioClient *minio.Client, bucketService *BucketService) *ObjectService {
	return &ObjectService{MinioClient: minioClient,
		BucketService: bucketService}
}

func (objectService *ObjectService) PutObject(ctx context.Context, bucket string, files []*multipart.FileHeader, folder string, tags ...map[string]string) ([]map[string]string, error) {
	if bucket == "" {
		return nil, errors.New("bucket can't be empty")

	}

	exists, err := objectService.BucketService.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("bucket doesn't exist")
	}

	var uploadInfoList []map[string]string
	var uuidFileName string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}

		defer src.Close()

		fileName := utils.GenerateUUIDFileName(file.Filename)

		if folder != "" {
			uuidFileName = folder + "/" + fileName

		} else {
			uuidFileName = fileName
		}

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
			"url": fmt.Sprintf("%s://%s/%s/%s/%s/%s/%s",
				ctx.Value("Scheme"),
				ctx.Value("Host"),
				"app/api/v1",
				"buckets",
				bucket,
				"files",
				uuidFileName,
			),
		})
	}

	return uploadInfoList, nil
}

func (objectService *ObjectService) GetObject(ctx context.Context, bucket, fileName string, options minio.GetObjectOptions) (*minio.Object, error) {
	var objectName string
	if bucket == "" || fileName == "" {
		return nil, errors.New("bucket or file name can't be empty")
	}

	exists, err := objectService.BucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("bucket doesn't exist")
	}

	if strings.Contains(fileName, "_") {
		folders := strings.Split(fileName, "_")
		objectName = folders[0] + "/" + folders[1]
	} else {
		objectName = fileName
	}
	return objectService.MinioClient.GetObject(ctx, bucket, objectName, options)
}

func (objectService *ObjectService) RemoveObjects(ctx context.Context, bucketName, objectName string, options minio.RemoveObjectOptions) error {
	err := objectService.MinioClient.RemoveObject(ctx, bucketName, objectName, options)
	return err
}

func (objectService *ObjectService) GetTag(ctx context.Context, bucket string, tagStr string) ([]map[string]string, error) {
	if bucket == "" {
		return nil, errors.New("bucket can't be empty")
	}

	exists, err := objectService.BucketService.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("bucket doesn't exist")
	}
	var uploadInfoList []map[string]string
	var existingTag = make(map[string]string)

	tags := make(map[string]string)
	if tagStr != "" {
		tagPairs := strings.Split(tagStr, ",")
		for _, tagPair := range tagPairs {
			pair := strings.Split(tagPair, "=")
			if len(pair) == 2 {
				tags[pair[0]] = pair[1]
			}
		}
	}

	// List objects in the bucket
	objectCh := objectService.MinioClient.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		tag, err := objectService.MinioClient.GetObjectTagging(ctx, bucket, object.Key, minio.GetObjectTaggingOptions{})
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
				"url": fmt.Sprintf("%s://%s/%s/%s/%s/%s/%s",
					ctx.Value("Scheme"),
					ctx.Value("Host"),
					"api/v1",
					"buckets",
					bucket,
					"files",
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

func (objectService *ObjectService) ObjectExists(existingObject []string, objectList []string) bool {
	for _, object := range objectList {
		found := false
		for _, existingObjects := range existingObject {
			if object == existingObjects {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (objectService *ObjectService) RemoveObjectTagging(background context.Context, bucket string, file string, options minio.RemoveObjectTaggingOptions) error {
	err := objectService.MinioClient.RemoveObjectTagging(background, bucket, file, options)
	return err
}
func (objectService *ObjectService) GetObjectTagging(background context.Context, bucket string, file string, options minio.GetObjectTaggingOptions) error {
	tags, err := objectService.MinioClient.GetObjectTagging(background, bucket, file, options)
	if tags == nil || err != nil {
		return err
	}
	return nil
}
