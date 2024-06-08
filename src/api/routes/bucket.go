package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	"cdn/src/redis"
	minio2 "cdn/src/service/minio"
	redis2 "cdn/src/service/redis"
	"github.com/gin-gonic/gin"
)

func BucketRoutes(router *gin.RouterGroup) {
	bucketService := minio2.NewBucketService(minio.GetInstance().GetMinio())
	objectService := minio2.NewObjectService(minio.GetInstance().GetMinio())
	redisService := redis2.NewRedisService(redis.GetInstance().GetClient())
	bucketController := controllers.NewBucketController(bucketService, objectService, redisService)

	storage := router.Group("storage")
	{
		// Bucket related routes
		storage.POST("buckets/:bucketName", bucketController.MakeBucket)
		storage.GET("buckets/:bucketName/objects", bucketController.ListObject)
		storage.DELETE("buckets/:bucketName", bucketController.RemoveBucket)
	}
}
