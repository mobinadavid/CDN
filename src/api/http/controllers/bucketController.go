package controllers

import (
	_ "cdn/src/api/http/requests"
	"cdn/src/api/http/response"
	i18n "cdn/src/pkg/i18h"
	"cdn/src/service/minio"
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

// MakeBucket handles make bucket requests
// @Summary Add new bucket
// @Description Adds a new bucket with the given details.
// @Tags Bucket
// @Param bucket query string  true "bucket"
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successMakeBucketRequest
// @Failure 400 {object} requests.failureMakeBucketRequest
// @Router /buckets/:bucket [post]
func (bucketController *BucketController) MakeBucket(c *gin.Context) {
	bucketName := c.Param("bucket")

	if bucketName == "" {
		response.Api(c).SetMessage("bucketName or region is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := bucketController.bucketService.BucketExists(c, bucketName)
	if exists || err != nil {
		response.Api(c).SetMessage("The specified bucket already exists.").SetStatusCode(http.StatusNotFound).Send()
		return
	}

	err = bucketController.bucketService.MakeBucket(c, bucketName, minio2.MakeBucketOptions{})
	if err != nil {
		response.Api(c).SetMessage("failed to create bucket.").SetStatusCode(http.StatusInternalServerError).Send()

		return
	}

	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"bucket": bucketName,
		}).Send()

}

// RemoveBucket handles bucket remove requests
// @Summary Delete bucket
// @Description Delete a bucket with the given uuid.
// @Param bucket query string true "bucket"
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveBucketRequest
// @Failure 400 {object} requests.failureRemoveBucketRequest
// @Router /buckets/:bucket [delete]
func (bucketController *BucketController) RemoveBucket(c *gin.Context) {
	bucketName := c.Param("bucket")

	exists, err := bucketController.bucketService.BucketExists(c, bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	// check if bucket is empty or not
	objects, err := bucketController.bucketService.ListObjects(c, bucketName, minio2.ListObjectsOptions{})
	if len(objects) != 0 {
		response.Api(c).SetMessage("The bucket is not empty.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	err = bucketController.bucketService.RemoveBucket(c, bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to remove bucket.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"bucket": bucketName,
		}).Send()

}

// ListObject handles pagination of objects.
// @Summary Get objects paginated data
// @Description Gets objects data with pagination.
// @Param bucket query string true "bucket"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetObjectListRequest
// @Failure 400 {object} requests.failureGetObjectListRequest
// @Router /buckets/:bucket/objects [get]
func (bucketController *BucketController) ListObject(c *gin.Context) {
	bucketName := c.Param("bucket")

	exists, err := bucketController.bucketService.BucketExists(c, bucketName)
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
			"info": object,
		})
	}

	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": objectList,
		}).Send()

}

// ListBucket handles pagination of buckets.
// @Summary Get buckets paginated data
// @Description Gets buckets data with pagination.
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetBucketListRequest
// @Failure 400 {object} requests.failureGetBucketListRequest
// @Router /buckets [get]
func (bucketController *BucketController) ListBucket(c *gin.Context) {
	buckets, err := bucketController.bucketService.ListBucket(c)
	if err != nil {
		response.Api(c).SetMessage("failed to list buckets.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"buckets": buckets,
		}).Send()
}
