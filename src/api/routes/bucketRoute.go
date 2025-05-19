package routes

import (
	"cdn/src/api/http/controllers"
	"cdn/src/minio"
	minio2 "cdn/src/service/minio"
	"github.com/gin-gonic/gin"
)

func BucketRoutes(router *gin.RouterGroup) {
	bucketService := minio2.NewBucketService(minio.GetInstance().GetInternalClient())
	objectService := minio2.NewObjectService(minio.GetInstance().GetInternalClient(),
		minio.GetInstance().GetCDNClient(), bucketService)
	bucketController := controllers.NewBucketController(bucketService, objectService)

	storage := router.Group("buckets")
	{
		// Bucket related routes
		storage.POST(":bucket", bucketController.MakeBucket)
		storage.GET(":bucket/objects", bucketController.ListObject)
		storage.DELETE(":bucket", bucketController.RemoveBucket)
		storage.GET("", bucketController.ListBucket)
	}
}
