package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	"cdn/src/service"
	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(router *gin.RouterGroup) {
	storageService := service.NewStorageService(minio.GetInstance().GetMinio(), "paresh")

	controller := controllers.NewStorageController(storageService)

	storage := router.Group("storage")
	{
		storage.POST("", controller.PutObject)
		//storage.GET(":fileName", controller.GetObject)
	}
}
