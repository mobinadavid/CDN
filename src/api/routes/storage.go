package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	"cdn/src/redis"
	"cdn/src/service"
	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(router *gin.RouterGroup) {
	storageService := service.NewStorageService(minio.GetInstance().GetMinio())
	redisService := service.NewRedisService(redis.GetInstance().GetClient())
	controller := controllers.NewStorageController(storageService, redisService)

	storage := router.Group("storage")
	{
		storage.POST("", controller.PutObject)
		storage.GET(":bucket/:file", controller.GetObject)
		storage.POST("buckets/:bucketName", controller.MakeBucket)
		storage.GET("buckets/:bucketName", controller.ListObject)
		storage.DELETE("buckets/:bucketName/", controller.RemoveObjects)
		storage.DELETE("buckets/:bucketName", controller.RemoveBucket)

	}
}
