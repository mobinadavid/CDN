package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/utils"
	"cdn/src/service"
	"context"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"strconv"
)

type objectController struct {
	bucketService *service.BucketService
	objectService *service.ObjectService
	redisService  *service.RedisService
}

func NewObjectController(bucketService *service.BucketService, objectService *service.ObjectService, redisService *service.RedisService) *objectController {
	return &objectController{bucketService: bucketService,
		objectService: objectService,
		redisService:  redisService,
	}
}
func (objectController *objectController) PutObject(c *gin.Context) {
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
	ip := c.ClientIP()
	//check rate limit

	allowed, err := objectController.redisService.CheckRateLimit(ip)
	if err != nil {
		response.Api(c).SetMessage("Failed to check rate limit").SetStatusCode(http.StatusInternalServerError).Send()
		return
	}
	if !allowed {
		response.Api(c).SetMessage("You can't put more than 10 Objects in one hour ").SetStatusCode(http.StatusTooManyRequests).Send()
		return
	}
	uploadInfoList, err := objectController.objectService.PutObject(context.WithValue(c.Request.Context(), "Scheme", c.GetHeader("Scheme")), bucket, form.File["files[]"])
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

func (objectController *objectController) GetObject(c *gin.Context) {
	bucketName := c.Param("bucketName")
	fileName := c.Param("file")

	if bucketName == "" || fileName == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	file, err := objectController.objectService.GetObject(context.Background(), bucketName, fileName, minio2.GetObjectOptions{})
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

func (objectController *objectController) RemoveObjects(c *gin.Context) {
	bucketName := c.Param("bucketName")
	exists, err := objectController.bucketService.BucketExists(context.Background(), bucketName)
	if err != nil {
		response.Api(c).SetMessage("failed to check if bucket exists.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	if !exists {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}
	objects, err := objectController.bucketService.ListObjects(context.Background(), bucketName, minio2.ListObjectsOptions{
		Recursive: true,
	})

	// Collect object names
	objectList := make([]string, 0, len(objects))
	for _, object := range objects {
		objectList = append(objectList, object.Key)
	}

	// Delete all objects

	for _, objectName := range objectList {
		errCh := objectController.objectService.RemoveObjects(context.Background(), bucketName, objectName, minio2.RemoveObjectOptions{})
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
