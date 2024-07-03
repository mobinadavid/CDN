package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/service/minio"
	"context"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"net/http"
)

type BucketController struct {
	bucketService *minio.BucketService
	objectService *minio.ObjectService
}

func NewBucketController(bucketService *minio.BucketService, objectService *minio.ObjectService) *BucketController {
	return &BucketController{bucketService: bucketService,
		objectService: objectService,
	}
}

func (bucketController *BucketController) MakeBucket(c *gin.Context) {

	bucketName := c.Param("bucket")

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
			"bucket": bucketName,
		}).Send()

}

func (bucketController *BucketController) RemoveBucket(c *gin.Context) {

	bucketName := c.Param("bucket")

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
			"bucket": bucketName,
		}).Send()

}

func (bucketController *BucketController) ListObject(c *gin.Context) {

	bucketName := c.Param("bucket")

	exists, err := bucketController.bucketService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	objects, err := bucketController.bucketService.ListObjects(c, bucketName, minio2.ListObjectsOptions{
		Recursive: true,
	})

	if err != nil {
		response.Api(c).SetMessage("failed to list objects.").SetStatusCode(http.StatusInternalServerError).Send()
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

func (bucketController *BucketController) ListBucket(c *gin.Context) {

	buckets, err := bucketController.bucketService.ListBucket(context.Background())
	if err != nil {
		response.Api(c).SetMessage("failed to list buckets.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	response.Api(c).
		SetMessage("Bucket created successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"buckets:": buckets,
		}).Send()

}
