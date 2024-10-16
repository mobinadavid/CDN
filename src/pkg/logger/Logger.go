package logger

import (
	"cdn/src/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"sync"
)

var (
	loggerInstance Logger
	once           sync.Once
)

type Logger interface {
	GetLogger() *zap.Logger
	SetLogger(logger *zap.Logger)
}

type Log struct {
	logger *zap.Logger
}

func Init() {
	GetInstance()
}

func GetInstance() Logger {
	once.Do(func() {
		core := zapcore.NewTee(zapcore.NewCore(encoderFile(), zapcore.NewMultiWriteSyncer(getLogger(), getRotator()), zapcore.DebugLevel))
		loggerInstance = &Log{
			logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		}
	})
	return loggerInstance
}

func encoderFile() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(conf)
}

func getRotator() zapcore.WriteSyncer {
	fileName := fmt.Sprintf("/log-data/%s", config.GetInstance().Get("LOG_FILE_NAME"))

	logSize, err := strconv.Atoi(config.GetInstance().Get("LOG_MAX_SIZE_MB"))
	if err != nil {
		logSize = 500
	}

	logMaxDay, err := strconv.Atoi(config.GetInstance().Get("LOG_MAX_DAY"))
	if err != nil {
		logMaxDay = 30
	}

	return zapcore.AddSync(
		&lumberjack.Logger{
			Filename: fileName,
			MaxSize:  logSize,
			MaxAge:   logMaxDay,
			Compress: true,
		})
}

func getLogger() zapcore.WriteSyncer {
	fileName := fmt.Sprintf("/log-data/%s", config.GetInstance().Get("LOG_FILE_NAME"))
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Return the rotator as the WriteSyncer
	return zapcore.AddSync(file)
}

func (l *Log) GetLogger() *zap.Logger {
	return l.logger
}

func (l *Log) SetLogger(logger *zap.Logger) {
	l.logger = logger
}
