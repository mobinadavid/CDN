package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"net/http"
)

type BucketController struct {
	bucketService *service.BucketService
	objectService *service.ObjectService
	redisService  *service.RedisService
}

func NewBucketController(bucketService *service.BucketService, objectService *service.ObjectService, redisService *service.RedisService) *BucketController {
	return &BucketController{bucketService: bucketService,
		objectService: objectService,
		redisService:  redisService,
	}
}
func (bucketController *BucketController) MakeBucket(c *gin.Context) {
	bucketName := c.Param("bucketName")

	if bucketName == "" {
		response.Api(c).SetMessage("bucketName or region is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	exists, err := bucketController.bucketService.BucketExists(context.Background(), bucketName)
	if exists || err != nil {
		response.Api(c).SetMessage("The specified bucket already exists.").SetStatusCode(http.StatusNotFound).Send()
		return
	}
	err = bucketController.bucketService.MakeBucket(context.Background(), bucketName, minio2.MakeBucketOptions{})
	if err != nil {
		response.Api(c).SetMessage("failed to create bucket.").SetStatusCode(http.StatusInternalServerError).Send()

		return
	}
	response.Api(c).
		SetMessage("Bucket created successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"bucketName": bucketName,
		}).Send()
}
func (bucketController *BucketController) RemoveBucket(c *gin.Context) {
	bucketName := c.Param("bucketName")
	exists, err := bucketController.bucketService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	// check if bucket is empty or not

	objects, err := bucketController.bucketService.ListObjects(context.Background(), bucketName, minio2.ListObjectsOptions{})
	if len(objects) != 0 {
		response.Api(c).SetMessage("The bucket is not empty.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	err = bucketController.bucketService.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to remove bucket.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}
	response.Api(c).
		SetMessage("Bucket is removed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"bucketName": bucketName,
		}).Send()

}
func (bucketController *BucketController) ListObject(c *gin.Context) {
	bucketName := c.Param("bucketName")
	exists, err := bucketController.bucketService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	objects, err := bucketController.bucketService.ListObjects(c, bucketName, minio2.ListObjectsOptions{})
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
			"number of objects:": len(objectList),
			"objects:":           objectList,
		}).Send()

}
