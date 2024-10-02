package logger

import (
	"cdn/src/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

var loggerInstance *zap.Logger
var fileName string

func Init() error {
	filePath := config.GetInstance().Get("LOG_FILE_PATH")
	if err := ensureDirExists(filePath); err != nil {
		return err
	}

	err := initLogger()
	if err != nil {
		return err
	}

	return nil
}

func GetInstance() *zap.Logger {
	newFileName := time.Now().Format("2006-01-02")
	if loggerInstance != nil && fileName == newFileName {
		return loggerInstance
	}

	_ = initLogger()
	fileName = newFileName
	return loggerInstance
}

func getCustomProductionLogger(filePath string) (*zap.Logger, error) {
	// New Production logger
	zapConfig := zap.NewProductionConfig()

	// Set the log level
	zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel())

	// Set file path
	zapConfig.OutputPaths = []string{filePath}
	zapConfig.ErrorOutputPaths = []string{filePath}

	// Build logger
	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func initLogger() error {
	// create a log file for each day
	fileDir := config.GetInstance().Get("LOG_FILE_PATH")
	filePath := fmt.Sprintf("%slog-%s.log", fileDir, time.Now().Format("2006-01-02"))

	logger, err := getCustomProductionLogger(filePath)
	if err != nil {
		return err
	}

	loggerInstance = logger
	return nil
}

// getLogLevel reads the log level from an environment variable or defaults to InfoLevel
func getLogLevel() zapcore.Level {
	level := config.GetInstance().Get("LOG_LEVEL")
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// ensureDirExists ensures that the directory for the log file exists
func ensureDirExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
