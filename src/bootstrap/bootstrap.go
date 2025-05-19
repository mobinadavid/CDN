package bootstrap

import (
	"cdn/src/api"
	"cdn/src/minio"
	"cdn/src/pkg/i18n"
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
	// Initialize Logger
	logger.Init()
	log.Println("Logger Service: Initialized Successfully.")

	// Initialize i18n
	err = i18n.Init()
	if err != nil {
		logger.GetInstance().GetLogger().Fatal("Failed to Initialize", zap.String("Service", "I18n"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	logger.GetInstance().GetLogger().Info("Initialized Successfully.", zap.String("Service", "I18n"), zap.Time("timestamp", time.Now()))

	//Initialize minio
	err = minio.Init()
	if err != nil {
		logger.GetInstance().GetLogger().Fatal("Failed to Initialize", zap.String("Service", "Minio"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	logger.GetInstance().GetLogger().Info("Initialized Successfully.", zap.String("Service", "Minio"), zap.Time("timestamp", time.Now()))

	//Initialize Redis
	err = redis.Init()
	if err != nil {
		logger.GetInstance().GetLogger().Fatal("Failed to Initialize", zap.String("Service", "Redis"), zap.Error(err), zap.Time("timestamp", time.Now()))
	}
	logger.GetInstance().GetLogger().Info("Initialized Successfully.", zap.String("Service", "Redis"), zap.Time("timestamp", time.Now()))

	// Initialize API
	go func() {
		err = api.Init()
		if err != nil {
			logger.GetInstance().GetLogger().Fatal("Failed to Initialize.", zap.String("Service", "API"), zap.Error(err), zap.Time("timestamp", time.Now()))
		}
		logger.GetInstance().GetLogger().Info("Initialized Successfully.", zap.String("Service", "API"), zap.Time("timestamp", time.Now()))
	}()
	logger.GetInstance().GetLogger().Info("Application is now running...", zap.Duration("duration_ms", time.Since(time.Now())), zap.Time("timestamp", time.Now()))
	logger.GetInstance().GetLogger().Info("Application is now running.\nPress CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logger.GetInstance().GetLogger().Info("Application is shutting down...", zap.Time("timestamp", time.Now()))
	log.Println("Application is shutting down...")

	return
}
