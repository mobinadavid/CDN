package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/api/http/middlewares"
	"cdn/src/minio"
	"cdn/src/redis"
	minio2 "cdn/src/service/minio"
	redis2 "cdn/src/service/redis"
	"github.com/gin-gonic/gin"
)

func ObjectRoutes(router *gin.RouterGroup) {
	bucketService := minio2.NewBucketService(minio.GetInstance().GetMinio())
	objectService := minio2.NewObjectService(minio.GetInstance().GetMinio())
	redisService := redis2.NewRedisService(redis.GetInstance().GetClient())
	objectController := controllers.NewObjectController(bucketService, objectService)

	storage := router.Group("storage")
	{
		storage.POST("", middlewares.RateLimit(redisService), objectController.PutObject)
		storage.GET("buckets/:bucketName/files/:file", objectController.GetObject)
		storage.DELETE("buckets/:bucketName/objects", objectController.RemoveObjects)

	}
}
