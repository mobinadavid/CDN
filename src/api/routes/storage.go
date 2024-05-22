package routes

import (
	"cdn/src/api/http/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(router *gin.RouterGroup) {
	controller := controllers.NewStorageController()

	storage := router.Group("storage")
	{
		storage.POST("", controller.PutObject)
		storage.GET(":bucket/:file", controller.GetObject)
	}
}
