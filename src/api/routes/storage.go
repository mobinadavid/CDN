package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	"cdn/src/redis"
	"cdn/src/service"
	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(router *gin.RouterGroup) {
	bucketService := service.NewBucketService(minio.GetInstance().GetMinio())
	objectService := service.NewObjectService(minio.GetInstance().GetMinio())
	redisService := service.NewRedisService(redis.GetInstance().GetClient())
	objectController := controllers.NewObjectController(bucketService, objectService, redisService)
	bucketController := controllers.NewBucketController(bucketService, objectService, redisService)

	storage := router.Group("storage")
	{
		storage.POST("", objectController.PutObject)
		storage.GET("buckets/:bucketName/files/:file", objectController.GetObject)
		storage.DELETE("buckets/:bucketName/objects", objectController.RemoveObjects)

		// Bucket related routes
		storage.POST("buckets/:bucketName", bucketController.MakeBucket)
		storage.GET("buckets/:bucketName/objects", bucketController.ListObject)
		storage.DELETE("buckets/:bucketName", bucketController.RemoveBucket)
	}
}
