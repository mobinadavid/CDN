package logger

import (
	"cdn/src/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

var loggerInstance *zap.Logger

func Init() error {
	return initLogger()
}

func GetInstance() *zap.Logger {
	if loggerInstance != nil {
		return loggerInstance
	}

	err := initLogger()
	if err != nil {
		return nil
	}
	return loggerInstance
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func getCustomProductionLogger() (*zap.Logger, error) {
	// get the file to put log
	logFile, err := getFile()
	if err != nil {
		return nil, err
	}

	// Create Rotator
	logSize, err := strconv.Atoi(config.GetInstance().Get("LOG_MAX_SIZE_MB"))
	if err != nil {
		logSize = 500
	}

	logMaxDay, err := strconv.Atoi(config.GetInstance().Get("LOG_MAX_DAY"))
	if err != nil {
		logMaxDay = 30
	}

	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename: logFile,
				MaxSize:  logSize,
				MaxAge:   logMaxDay,
				Compress: true,
			},
		}, nil
	})

	// New Production logger
	zapConfig := zap.NewProductionConfig()

	// Set the log level
	zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel())

	// Set file path
	zapConfig.OutputPaths = []string{fmt.Sprintf("lumberjack:%s", logFile)}
	zapConfig.ErrorOutputPaths = []string{logFile, "stderr"}

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

func getLocalLogger() (*zap.Logger, error) {
	// get the file to put log
	logFile, err := getFile()
	if err != nil {
		return nil, err
	}

	// New Production logger
	zapConfig := zap.NewProductionConfig()

	// Set the log level
	zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel())

	// Set file path
	zapConfig.OutputPaths = []string{logFile, "stdout"}
	zapConfig.ErrorOutputPaths = []string{logFile, "stderr"}

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
	appEnv := config.GetInstance().Get("APP_ENV")
	if appEnv == "production" {
		logger, err := getCustomProductionLogger()
		if err != nil {
			return err
		}
		loggerInstance = logger
	} else if appEnv == "local" || appEnv == "dev" || appEnv == "development" {
		logger, err := getLocalLogger()
		if err != nil {
			return err
		}
		loggerInstance = logger
	}
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

func getFile() (string, error) {
	filePath := config.GetInstance().Get("LOG_DIR_NAME")
	err := ensureDirExists(filePath)
	if err != nil {
		return "", err
	}

	fileName := config.GetInstance().Get("LOG_FILE_NAME")

	if config.GetInstance().Get("APP_ENV") == "production" {
		return fmt.Sprintf("/%s/%s", filePath, fileName), nil
	}
	return fmt.Sprintf("%s/%s", filePath, fileName), nil
}
