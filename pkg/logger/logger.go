package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() *zap.Logger {
	logLevelStr := env.GetLoggingLevelEnv("LOG_LEVEL", "INFO")
	var logLevel zapcore.Level
	
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		logLevel = zapcore.DebugLevel
	case "INFO":
		logLevel = zapcore.InfoLevel
	case "WARNING":
		logLevel = zapcore.WarnLevel
	case "ERROR":
		logLevel = zapcore.ErrorLevel
	default:
		fmt.Printf("Неизвестный уровень логирования: %s, установлен INFO\n", logLevelStr)
		logLevel = zapcore.InfoLevel
	}

	// Создаём конфигурацию логгера
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	config.Encoding = "console"

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
		os.Exit(1) // Завершаем программу
	}

	return logger
}