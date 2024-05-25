package controllers

import (
	"cdn/src/api/http/response"
	"cdn/src/pkg/utils"
	"cdn/src/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

//var bucketName = "paresh" //todo: hard-coded

type StorageController struct {
	storageService *service.StorageService
}

func NewStorageController(storageService *service.StorageService) *StorageController {
	return &StorageController{
		storageService: storageService,
	}
}

func (storageController *StorageController) PutObject(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	files := form.File["files[]"]
	err = utils.ValidateFiles(files)
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnprocessableEntity).Send()
		return
	}

	var uploadInfoList []map[string]string
	for _, file := range files {
		uploadInfo, err := storageController.storageService.UploadFile(c, file)
		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			return
		}
		uploadInfoList = append(uploadInfoList, uploadInfo)
	}

	response.Api(c).
		SetMessage("Files uploaded successfully").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"objects": uploadInfoList,
		}).Send()
}

func (storageController *StorageController) GetObject(c *gin.Context) {
	fileName := c.Param("fileName")
	file, err := storageController.storageService.GetObject(fileName)
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
