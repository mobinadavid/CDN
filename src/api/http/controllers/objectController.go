package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/utils"
	"cdn/src/service/minio"
	"context"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"strconv"
	"strings"
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
// @Tags CDN
// @Param files formData file true "File to upload"
// @Param bucket formData string true "Bucket name"
// @Success 200 {object} requests.successPutObjectRequest
// @Failure 400 {object} requests.failurePutObjectRequest
// @Router /storage [post]
func (objectController *ObjectController) PutObject(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	bucket := c.PostForm("bucket")
	folder := c.PostForm("folder")
	tagsStr := c.PostForm("tag")

	if bucket == "" {
		response.Api(c).SetMessage("bucket is required.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return

	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)

	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if err := utils.ValidateFiles(form.File["files[]"]); err != nil {
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

// GetObject handles get data of object.
// @Summary Get object
// @Description Gets object data with specified filename.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags CDN
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetObjectRequest
// @Failure 400 {object} requests.failureGetObjectRequest
// @Router /storage/buckets/:bucket/files/:file [get]
func (objectController *ObjectController) GetObject(c *gin.Context) {

	var objectName string

	bucket := c.Param("bucket")
	fileName := c.Param("file")

	if bucket == "" || fileName == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if strings.Contains(fileName, "_") {
		folders := strings.Split(fileName, "_")
		objectName = folders[0] + "/" + folders[1]
	} else {
		objectName = fileName
	}

	file, err := objectController.objectService.GetObject(context.Background(), bucket, objectName, minio2.GetObjectOptions{})
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

// RemoveObjects handles objects remove requests
// @Summary Delete objects of a bucket
// @Description Delete objects of a bucket.
// @Param bucket query string true "bucket"
// @Tags CDN
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveObjectsRequest
// @Failure 400 {object} requests.failureRemoveObjectsRequest
// @Router /storage/buckets/:bucket/objects [delete]
func (objectController *ObjectController) RemoveObjects(c *gin.Context) {

	bucket := c.Param("bucket")

	if bucket == "" {
		response.Api(c).SetMessage("bucket is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	objects, err := objectController.bucketService.ListObjects(context.Background(), bucket, minio2.ListObjectsOptions{
		Recursive: true,
	})

	// Collect object names
	objectList := make([]string, 0, len(objects))
	for _, object := range objects {
		objectList = append(objectList, object.Key)
	}

	// Delete all objects
	for _, objectName := range objectList {
		errCh := objectController.objectService.RemoveObjects(context.Background(), bucket, objectName, minio2.RemoveObjectOptions{})
		if errCh != nil {
			response.Api(c).SetMessage("failed to remove objects.").SetStatusCode(http.StatusInternalServerError).Send()
			return
		}
	}

	response.Api(c).
		SetMessage("all removed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"object's name:": objectList,
		}).Send()

}

// RemoveObject handles object remove requests
// @Summary Delete object
// @Description Delete an object with the file.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags CDN
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveObjectRequest
// @Failure 400 {object} requests.failureRemoveObjectRequest
// @Router /storage/buckets/:bucket/files/:file [delete]
func (objectController *ObjectController) RemoveObject(c *gin.Context) {

	bucket := c.Param("bucket")
	fileName := c.Param("file")

	if bucket == "" || fileName == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	objects, err := objectController.bucketService.ListObjects(context.Background(), bucket, minio2.ListObjectsOptions{
		Recursive: true,
	})

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

	exist := objectController.objectService.ObjectExists(existingObjectList, objectList)
	if exist == false {
		response.Api(c).SetMessage("failed to find objects.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	// Delete the objects
	for _, objectName := range objectList {

		errCh := objectController.objectService.RemoveObjects(context.Background(), bucket, objectName, minio2.RemoveObjectOptions{})

		if errCh != nil {
			response.Api(c).SetMessage("failed to remove objects.").SetStatusCode(http.StatusInternalServerError).Send()
			return
		}

	}

	response.Api(c).
		SetMessage("removed successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"object's name:": objectList,
		}).Send()

}

// GetTag handles get data of tag.
// @Summary Get tag
// @Description Gets tag data.
// @Param bucket query string true "bucket"
// @Param file query string true "file"
// @Tags CDN
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successGetTagRequest
// @Failure 400 {object} requests.failureGetTagRequest
// @Router /storage/buckets/:bucket/tags/:tag [get]
func (objectController *ObjectController) GetTag(c *gin.Context) {

	bucket := c.Param("bucket")
	tagsStr := c.Param("tag")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Host", c.Request.Host)
	ctx = context.WithValue(ctx, "Scheme", c.GetHeader("Scheme"))

	if bucket == "" {
		response.Api(c).SetMessage("bucket is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	urls, err := objectController.objectService.GetTag(ctx, bucket, tagsStr)
	if err != nil {
		response.Api(c).SetMessage("Failed to get objects by tags.").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	response.Api(c).
		SetMessage("urls retrieved successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": urls,
		}).Send()

}

// RemoveTag handles tag remove requests
// @Summary Delete tag
// @Description Delete a tag with the given object.
// @Param bucket query string true "bucket"
// @Param object query string true "object"
// @Tags CDN
// @Accept  json
// @Produce  json
// @Success 200 {object} requests.successRemoveTagRequest
// @Failure 400 {object} requests.failureRemoveTagRequest
// @Router /storage/buckets/:bucket/objects/:object [delete]
func (objectController *ObjectController) RemoveTag(c *gin.Context) {

	bucket := c.Param("bucket")
	file := c.Param("object")

	if bucket == "" || file == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	exists, err := objectController.bucketService.BucketExists(context.Background(), bucket)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
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
		response.Api(c).SetMessage("object doesnt have any tag").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}

	err = objectController.objectService.RemoveObjectTagging(context.Background(), bucket, file, minio2.RemoveObjectTaggingOptions{})
	if err != nil {
		response.Api(c).SetMessage("can't remove the tag").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	response.Api(c).
		SetMessage("removed successfully").
		SetStatusCode(http.StatusOK).Send()
}
