package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	once     sync.Once
	instance *zap.Logger
)

func initializeLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Enable color output
	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	instance = logger
}

func HireScribe() *zap.Logger {
	once.Do(initializeLogger)
	return instance
}
