package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	"cdn/src/service"
	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(router *gin.RouterGroup) {
	service := service.NewStorageService(minio.GetInstance().GetMinio())
	controller := controllers.NewStorageController(service)

	storage := router.Group("storage")
	{
		storage.POST("", controller.PutObject)
		storage.GET(":bucket/:file", controller.GetObject)
		storage.POST(":bucketName/:region", controller.MakeBucket)
	}
}
