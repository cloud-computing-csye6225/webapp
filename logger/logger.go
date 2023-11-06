package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"sync"
)

var once sync.Once

var logger *zap.Logger

func InitLogger() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "D:/CSYE6225/app/logs/app.log",
			MaxSize:    1,
			MaxBackups: 10,
			MaxAge:     14,
			Compress:   false,
		})

		level := zap.InfoLevel
		levelEnv := os.Getenv("LOG_LEVEL")
		if levelEnv != "" {
			levelFromEnv, err := zapcore.ParseLevel(levelEnv)
			if err != nil {
				log.Println(
					fmt.Errorf("invalid level, defaulting to INFO: %w", err),
				)
			}
			level = levelFromEnv
		}

		logLevel := zap.NewAtomicLevelAt(level)

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleEncoder := zapcore.NewJSONEncoder(productionCfg)
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel),
			zapcore.NewCore(fileEncoder, file, logLevel),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})

	return logger
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
func Debug(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
func Error(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
func Fatal(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
