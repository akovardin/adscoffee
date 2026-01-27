package logger

import (
	"go.uber.org/zap"
)

func New() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
