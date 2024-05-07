package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/minio"
	"fmt"
	"github.com/gin-gonic/gin"
	minio2 "github.com/minio/minio-go/v7"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var bucketName = "paresh" //todo: hard-coded

type UploadInfo struct {
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
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
	form, _ := c.MultipartForm()

	err := validateFiles(form.File["files"])

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

		_, err = storageController.minio.GetMinio().PutObject(c, bucketName, file.Filename, src, file.Size, minio2.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}

		presignedURL, err := storageController.minio.GetMinio().PresignedGetObject(c, bucketName, file.Filename, 8*time.Hour, nil)

		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}

		uploadInfoList = append(uploadInfoList, UploadInfo{
			FileName: file.Filename,
			Size:     file.Size,
			URL:      presignedURL.String(),
		})
	}

	response.Api(c).SetMessage("Files uploaded successfully").SetStatusCode(http.StatusOK).SetData(map[string]interface{}{
		"uploads": uploadInfoList,
	}).Send()

	return

}

// Custom validation function for uploaded files
func validateFiles(files []*multipart.FileHeader) error {
	// Allowed file extensions todo: hard-coded
	allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf", "zip", "rar", "docx", "doc", "csv", "xlsx", "mkv", "mp4"}

	for _, file := range files {
		ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))
		// Check if file extension is allowed
		valid := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid file format: %s", file.Filename)
		}
	}

	return nil
}
