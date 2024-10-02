package bootstrap

import (
	"cdn/src/api"
	"cdn/src/minio"
	"cdn/src/pkg/logger"
	"cdn/src/redis"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init() (err error) {
	defer func() {
		log.Println("Goodbye!")
		os.Exit(0)
	}()

	//Initialize logger
	err = logger.Init()
	if err != nil {
		log.Fatalf("Logger Service: Failed to Initialize. %v", err)
	}
	logger.GetInstance().Info("Initialized Successfully.", zap.String("Service", "Logger"), zap.Time("timestamp", time.Now()))
	log.Println("Logger Service: Initialized Successfully.")

	//Initialize minio
	err = minio.Init()
	if err != nil {
		log.Println("Minio Service: Failed to Initialize. %v", err)
		logger.GetInstance().Fatal("Failed to Initialize", zap.String("Service", "Minio"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	log.Println("Minio Service: Initialized Successfully.")

	//Initialize Redis
	err = redis.Init()
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		logger.GetInstance().Fatal("Failed to Initialize", zap.String("Service", "Redis"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	logger.GetInstance().Info("Initialized Successfully.", zap.String("Service", "Redis"), zap.Time("timestamp", time.Now()))
	log.Println("Redis Service: Initialized Successfully.")

	// Initialize API
	go func() {
		err = api.Init()
		if err != nil {
			log.Println("API Service: Failed to Initialize. %v", err)
			logger.GetInstance().Info("Failed to Initialize.", zap.String("Service", "API"), zap.Error(err), zap.Time("timestamp", time.Now()))
		}
		logger.GetInstance().Info("Initialized Successfully.", zap.String("Service", "API"), zap.Time("timestamp", time.Now()))
		log.Println("API Service: Initialized Successfully.")
	}()
	logger.GetInstance().Info("Application is now running...", zap.Duration("duration_ms", time.Since(time.Now())), zap.Time("timestamp", time.Now()))
	log.Println("Application is now running.\nPress CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logger.GetInstance().Info("Application is shutting down...", zap.Time("timestamp", time.Now()))
	log.Println("Application is shutting down...")

	return
}
