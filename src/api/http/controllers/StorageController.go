package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/minio"
	"cdn/src/pkg/utils"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"net/http"
	"time"
)

var bucketName = "paresh" //todo: hard-coded

type UploadInfo struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
}

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

	err = utils.ValidateFiles(form.File["files"])

	if err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	var uploadInfoList []UploadInfo

	for _, file := range form.File["files[]"] {
		src, err := file.Open()
		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}
		defer src.Close()

		uuidFileName := utils.GenerateUUIDFileName(file.Filename)

		_, err = storageController.minio.GetMinio().PutObject(c, bucketName, uuidFileName, src, file.Size, minio2.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

		if err != nil {
			response.Api(c).
				SetMessage(err.Error()).
				SetStatusCode(http.StatusInternalServerError).
				Send()

			return
		}

		preSignedURL, err := storageController.minio.GetMinio().PresignedGetObject(c, bucketName, uuidFileName, (7*24)*time.Hour, nil)

		if err != nil {
			response.Api(c).
				SetMessage(err.Error()).
				SetStatusCode(http.StatusInternalServerError).
				Send()

			return
		}

		uploadInfoList = append(uploadInfoList, UploadInfo{
			FileName: uuidFileName,
			URL:      preSignedURL.String(),
		})
	}

	response.Api(c).
		SetMessage("Files uploaded successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"uploads": uploadInfoList,
		}).Send()

	return

}
