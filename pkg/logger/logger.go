package logger

import (
	"go.uber.org/zap"
)

// Logger carries the loggers
type Logger struct {
	Zap *zap.Logger
}

// New creates a global logger
func New() *Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{Zap: zapLogger}
}
