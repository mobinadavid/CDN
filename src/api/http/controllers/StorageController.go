package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/minio"
	"cdn/src/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type StorageController struct {
	minio *minio.Client
}

func NewStorageController() *StorageController {
	return &StorageController{
		minio: minio.GetInstance(),
	}
}

func (storageController *StorageController) PutObject(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	if c.PostForm("bucket") == "" {
		response.Api(c).SetMessage("bucket is required.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	bucketExists, err := storageController.minio.GetMinio().BucketExists(context.Background(), c.PostForm("bucket"))
	if !bucketExists || err != nil {
		response.Api(c).SetMessage("The specified bucket does not exist.").SetStatusCode(http.StatusUnprocessableEntity).Send()
	}

	err = utils.ValidateFiles(form.File["files[]"])
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	var uploadInfoList []map[string]string

	for _, file := range form.File["files[]"] {
		src, err := file.Open()
		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}
		defer src.Close()

		uuidFileName := utils.GenerateUUIDFileName(file.Filename)

		_, err = storageController.minio.GetMinio().PutObject(c, c.PostForm("bucket"), uuidFileName, src, file.Size, minio2.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

		if err != nil {
			response.Api(c).
				SetMessage(err.Error()).
				SetStatusCode(http.StatusInternalServerError).
				Send()

			return
		}

		uploadInfoList = append(uploadInfoList, map[string]string{
			"original_file_name": strings.ToLower(file.Filename),
			"size":               strconv.FormatInt(file.Size, 10),
			"file_name":          uuidFileName,
			"url": fmt.Sprintf("%s://%s/%s/%s/%s",
				c.GetHeader("Scheme"),
				c.Request.Host,
				"app/api/v1/storage",
				c.PostForm("bucket"),
				uuidFileName,
			),
		})
	}

	response.Api(c).
		SetMessage("Files uploaded successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": uploadInfoList,
		}).Send()
	return
}

func (storageController *StorageController) GetObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	fileName := c.Param("file")

	if bucketName == "" || fileName == "" {
		response.Api(c).SetMessage("bucket or file is missing.").SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	file, err := storageController.minio.GetMinio().GetObject(context.Background(), bucketName, fileName, minio2.GetObjectOptions{})
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
