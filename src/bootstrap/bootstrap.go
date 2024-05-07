package bootstrap

import (
	"cdn/src/api"
	"cdn/src/pkg/minio"
	"context"
	minio2 "github.com/minio/minio-go/v7"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var bucketName = "paresh" //todo: hard-coded

func Init() (err error) {

	defer func() {
		log.Println("Goodbye!")
		os.Exit(0)
	}()

	// Initialize Minio
	err = minio.Init()
	if err != nil {
		log.Fatalf("Minio Service: Failed to Initialize. %v", err)
	}
	log.Println("Minio Service: Initialized Successfully.")

	// Create Bucket todo:Hard-coded

	minioInstance := minio.GetInstance().GetMinio()
	bucketExists, err := minioInstance.BucketExists(context.Background(), bucketName)

	if !bucketExists {
		err = minio.GetInstance().GetMinio().MakeBucket(context.Background(), bucketName, minio2.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Minio Service: Failed to Create Bucket. %v", err)
		}
	}

	// Initialize API
	go func() {
		err = api.Init()
		if err != nil {
			log.Fatalf("API Service: Failed to Initialize. %v", err)
		}
		log.Println("API Service: Initialized Successfully.")
	}()

	log.Println("Application is now running.\nPress CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Application is shutting down...")

	return
}
