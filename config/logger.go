package config

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			log.Fatalf("❌ Failed to create logs directory: %v", err)
		}
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "logs/api.log"}

	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatalf("❌ Cannot initialize logger: %v", err)
	}
	Logger = zapLogger.Sugar()
}
