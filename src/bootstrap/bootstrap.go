package bootstrap

import (
	"cdn/src/api"
	"cdn/src/minio"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
	//Initialize Redis
	//err = drivers.Init()
	//if err != nil {
	//	log.Fatalln("Failed to connect to Redis:", err)
	//}
	//log.Println("Redis Service: Initialized Successfully.")
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
