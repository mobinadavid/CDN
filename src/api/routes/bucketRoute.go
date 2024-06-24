package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	minio2 "cdn/src/service/minio"
	"github.com/gin-gonic/gin"
)

func BucketRoutes(router *gin.RouterGroup) {
	bucketService := minio2.NewBucketService(minio.GetInstance().GetMinio())
	objectService := minio2.NewObjectService(minio.GetInstance().GetMinio())
	bucketController := controllers.NewBucketController(bucketService, objectService)

	storage := router.Group("storage")
	{
		// Bucket related routes
		storage.POST("buckets/:bucket", bucketController.MakeBucket)
		storage.GET("buckets/:bucket/objects", bucketController.ListObject)
		storage.DELETE("buckets/:bucket", bucketController.RemoveBucket)
	}
}
