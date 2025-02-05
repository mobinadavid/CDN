package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/config"
	"cdn/src/pkg/i18n"
	"cdn/src/pkg/logger"
	"cdn/src/pkg/utils"
	"cdn/src/service/minio"
	"context"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ObjectController struct {
	bucketService *minio.BucketService
	objectService *minio.ObjectService
}

func NewObjectController(bucketService *minio.BucketService, objectService *minio.ObjectService) *ObjectController {
	return &ObjectController{bucketService: bucketService,
		objectService: objectService,
	}
}

// PutObject handles put object on bucket requests
// @Summary Add new object to bucket
// @Description Adds a new object to bucket with the given details.
// @Tags Object
// @Accept multipart/form-data
// @Param files formData file true "File to upload"
// @Param bucket formData string true "Bucket name"
// @Success 200 {object} requests.successPutObjectRequest
// @Failure 400 {object} requests.failurePutObjectRequest
// @Router /buckets [post]
func (objectController *ObjectController) PutObject(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Put Object"))
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	bucket := c.PostForm("bucket")
	folder := c.PostForm("folder")
	tagsStr := c.PostForm("tag")

	if err := utils.ValidateFiles(form.File["files[]"]); err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Put Object"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	tags := make(map[string]string)
	var uploadInfoList []map[string]string

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Host", c.Request.Host)
	ctx = context.WithValue(ctx, "Scheme", c.GetHeader("Scheme"))

	switch {
	case tagsStr != "":
		tagPairs := strings.Split(tagsStr, ",")
		for _, tagPair := range tagPairs {
			pair := strings.Split(tagPair, "=")
			if len(pair) == 2 {
				tags[pair[0]] = pair[1]
			}
		}

		uploadInfoList, err = objectController.objectService.PutObject(ctx, bucket, form.File["files[]"], folder, tags)
	default:
		uploadInfoList, err = objectController.objectService.PutObject(ctx, bucket, form.File["files[]"], folder)
	}

	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Put Object"))
		response.Api(c).
			SetMessage(err.Error()).
			SetStatusCode(http.StatusInternalServerError).
			Send()
		return
	}

	// Return response.
	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": uploadInfoList,
		}).
		Send()
}

// GetObject handles get data of object.
// @Summary Get object
// @Description Gets object data with specified filename.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetObjectRequest
// @Failure 400 {object} requests.failureGetObjectRequest
// @Router /buckets/:bucket/files/:file [get]
func (objectController *ObjectController) GetObject(c *gin.Context) {
	bucket := c.Param("bucket")
	fileName := c.Param("file")

	file, err := objectController.objectService.GetObject(c, bucket, fileName, minio2.GetObjectOptions{})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Object"))
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Object"))
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))

	_, err = io.Copy(c.Writer, file)

	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Object"))
		response.Api(c).SetStatusCode(http.StatusInternalServerError).Send()
		return
	}
}

// GetPreSigned handles get preSigned url of object.
// @Summary Get PreSigned
// @Description Gets object data with specified filename.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetObjectRequest
// @Failure 400 {object} requests.failureGetObjectRequest
// @Router /buckets/:bucket/files/url/:file [get]
func (objectController *ObjectController) GetPreSigned(c *gin.Context) {
	bucket := c.Param("bucket")
	fileName := c.Param("file")

	file, err := objectController.objectService.GetObject(c, bucket, fileName, minio2.GetObjectOptions{})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Presigned"))
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}
	defer file.Close()

	expiryNum := config.GetInstance().Get("MINIO_PRE_SIGNED_URL_EXPIRE_TIME")
	expiry, err := strconv.Atoi(expiryNum)
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Presigned"))
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	preSignedURL, err := objectController.objectService.ClientCdn.PresignedGetObject(ctx, bucket, fileName, time.Duration(expiry)*time.Minute, nil)
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Presigned"))
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	// Return response.
	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"file_name": fileName,
			"url":       preSignedURL.String(),
		},
		).Send()
}

// RemoveObjects handles objects remove requests
// @Summary Delete objects of a bucket
// @Description Delete objects of a bucket.
// @Param bucket query string true "bucket"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveObjectsRequest
// @Failure 400 {object} requests.failureRemoveObjectsRequest
// @Router /buckets/:bucket/objects [delete]
func (objectController *ObjectController) RemoveObjects(c *gin.Context) {
	bucket := c.Param("bucket")

	objects, err := objectController.bucketService.ListObjects(context.Background(), bucket, minio2.ListObjectsOptions{
		Recursive: true,
	})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Remove Objects"))
		response.Api(c).SetStatusCode(http.StatusInternalServerError).Send()
	}

	// Collect object names
	objectList := make([]string, 0, len(objects))
	for _, object := range objects {
		objectList = append(objectList, object.Key)
	}

	// Delete all objects
	for _, objectName := range objectList {
		errCh := objectController.objectService.RemoveObjects(context.Background(), bucket, objectName, minio2.RemoveObjectOptions{})
		if errCh != nil {
			logger.GetInstance().Error(errCh.Error(), zap.String("Method", "Remove Objects"))
			response.Api(c).SetMessage(errCh.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}
	}

	// Return response.
	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": objectList,
		}).
		Send()
}

// RemoveObject handles object remove requests
// @Summary Delete object
// @Description Delete an object with the file.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveObjectRequest
// @Failure 400 {object} requests.failureRemoveObjectRequest
// @Router /buckets/:bucket/files/:file [delete]
func (objectController *ObjectController) RemoveObject(c *gin.Context) {
	bucket := c.Param("bucket")
	fileName := c.Param("file")

	if fileName == "" {
		response.Api(c).SetMessage("file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	objects, err := objectController.bucketService.ListObjects(c, bucket, minio2.ListObjectsOptions{
		Recursive: true,
	})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Remove Object"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
	}

	// Collect object names
	existingObjectList := make([]string, 0, len(objects))
	for _, object := range objects {
		existingObjectList = append(existingObjectList, object.Key)
	}

	objectList := make([]string, 0)
	for _, objectName := range strings.Split(fileName, ",") {
		if strings.Contains(objectName, "_") {
			splitString := strings.Split(objectName, "_")
			objectName = splitString[0] + "/" + splitString[1]
		}
		objectList = append(objectList, objectName)
	}

	// Delete the objects
	for _, objectName := range objectList {
		errCh := objectController.objectService.RemoveObjects(c, bucket, objectName, minio2.RemoveObjectOptions{})
		if errCh != nil {
			logger.GetInstance().Error(errCh.Error(), zap.String("Method", "Remove Object"))
			response.Api(c).SetMessage(errCh.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}

	}

	// Return response.
	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": objectList,
		}).
		Send()
}

// GetTag handles get data of tag.
// @Summary Get tag
// @Description Gets tag data.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetTagRequest
// @Failure 400 {object} requests.failureGetTagRequest
// @Router /buckets/:bucket/tags/:tag [get]
func (objectController *ObjectController) GetTag(c *gin.Context) {
	bucket := c.Param("bucket")
	tagsStr := c.Param("tag")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Host", c.Request.Host)
	ctx = context.WithValue(ctx, "Scheme", c.GetHeader("Scheme"))

	names, err := objectController.objectService.GetTag(ctx, bucket, tagsStr)
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Get Tag"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	// Return response.
	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": names,
		}).
		Send()
}

// RemoveTag handles tag remove requests
// @Summary Delete tag
// @Description Delete a tag with the given object.
// @Param bucket query string true "bucket"
// @Param object query string true "object"
// @Tags Object
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveTagRequest
// @Failure 400 {object} requests.failureRemoveTagRequest
// @Router /buckets/:bucket/objects/:object [delete]
func (objectController *ObjectController) RemoveTag(c *gin.Context) {
	bucket := c.Param("bucket")
	file := c.Param("object")

	if bucket == "" || file == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Remove Tag"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if strings.Contains(file, "_") {
		splitString := strings.Split(file, "_")
		file = splitString[0] + "/" + splitString[1]
	}

	err = objectController.objectService.GetObjectTagging(context.Background(), bucket, file, minio2.GetObjectTaggingOptions{})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Remove Tag"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	err = objectController.objectService.RemoveObjectTagging(context.Background(), bucket, file, minio2.RemoveObjectTaggingOptions{})
	if err != nil {
		logger.GetInstance().Error(err.Error(), zap.String("Method", "Remove Tag"))
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	response.Api(c).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetStatusCode(http.StatusOK).
		Send()
}
