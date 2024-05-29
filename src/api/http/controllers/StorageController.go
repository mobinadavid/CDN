package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/utils"
	"cdn/src/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"strconv"
)

type StorageController struct {
	storageService *service.StorageService
}

func NewStorageController(storageService *service.StorageService) *StorageController {
	return &StorageController{storageService: storageService}
}

func (storageController *StorageController) PutObject(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	bucket := c.PostForm("bucket")
	if bucket == "" {
		response.Api(c).SetMessage("bucket is required.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return

	}

	exists, err := storageController.storageService.BucketExists(context.Background(), bucket)
	if !exists || err != nil {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if err := utils.ValidateFiles(form.File["files[]"]); err != nil {
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	uploadInfoList, err := storageController.storageService.UploadFiles(context.WithValue(c.Request.Context(), "Scheme", c.GetHeader("Scheme")), bucket, form.File["files[]"])
	if err != nil {
		response.Api(c).
			SetMessage(err.Error()).
			SetStatusCode(http.StatusInternalServerError).
			Send()
		return
	}

	response.Api(c).
		SetMessage("Files uploaded successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": uploadInfoList,
		}).Send()
}

func (storageController *StorageController) GetObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	fileName := c.Param("file")

	if bucketName == "" || fileName == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	file, err := storageController.storageService.GetObject(context.Background(), bucketName, fileName, minio2.GetObjectOptions{})
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))

	_, err = io.Copy(c.Writer, file)

	if err != nil {
		response.Api(c).SetStatusCode(http.StatusInternalServerError).Send()
		return
	}
}
func (storageController *StorageController) MakeBucket(c *gin.Context) {
	bucketName := c.Param("bucketName")
	//region is hardCoded
	//us-east-1
	region := c.Param("region")

	if bucketName == "" || region == "" {
		response.Api(c).SetMessage("bucketName or region is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	exists, err := storageController.storageService.BucketExists(context.Background(), bucketName)
	if exists || err != nil {
		response.Api(c).SetMessage("The specified bucket already exists.").SetStatusCode(http.StatusNotFound).Send()
		return
	}
	err = storageController.storageService.MakeBucket(context.Background(), bucketName, minio2.MakeBucketOptions{Region: region})
	if err != nil {
		response.Api(c).SetMessage("failed to create bucket.").SetStatusCode(http.StatusInternalServerError).Send()

		return
	}
	response.Api(c).
		SetMessage("Bucket created successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"name":     bucketName,
			"location": region,
		}).Send()
}
func (storageController *StorageController) RemoveBucket(c *gin.Context) {
	bucketName := c.Param("bucketName")
	exists, err := storageController.storageService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	err = storageController.storageService.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to remove bucket.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}
	response.Api(c).
		SetMessage("Bucket removed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"name": bucketName,
		}).Send()

}
func (storageController *StorageController) ListObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	exists, err := storageController.storageService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	objects, err := storageController.storageService.ListObjects(c, bucketName, minio2.ListObjectsOptions{})
	if err != nil {
		response.Api(c).SetMessage("failed to list objects.").SetStatusCode(http.StatusInternalServerError).Send()
		fmt.Println(err)
		return
	}
	objectList := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		objectList = append(objectList, map[string]interface{}{
			"info:": object,
		})
	}

	response.Api(c).
		SetMessage("listed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects:": objectList,
		}).Send()

}
func (storageController *StorageController) RemoveObjects(c *gin.Context) {
	bucketName := c.Param("bucketName")
	exists, err := storageController.storageService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	objects, err := storageController.storageService.ListObjects(context.Background(), bucketName, minio2.ListObjectsOptions{
		Recursive: true,
	})

	// Collect object names
	objectList := make([]string, 0, len(objects))
	for _, object := range objects {
		objectList = append(objectList, object.Key)
	}

	// Delete all objects

	for _, objectName := range objectList {
		errCh := storageController.storageService.RemoveObjects(context.Background(), bucketName, objectName, minio2.RemoveObjectOptions{})
		{
			if errCh != nil {
				response.Api(c).SetMessage("failed to remove objects.").SetStatusCode(http.StatusInternalServerError).Send()
				return
			}

		}

	}
	response.Api(c).
		SetMessage("all removed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"object's name:": objectList,
		}).Send()

}
