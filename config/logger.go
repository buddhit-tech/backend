package config

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger = nil

func InitLogger() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap_logger_init: error=%s", err.Error())
	}

	logger = zapLogger
}

func GetLogger() *zap.Logger {
	return logger
}
