package utils

import (
	"context"

	"go.uber.org/zap"
)

func GetLogger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(LoggerKey).(*zap.Logger)
	return logger
}

func CreateContextWithLogger(ctx context.Context, loggerName string) context.Context {
	baseLogger := GetLogger(ctx)

	newLogger := baseLogger.Named(loggerName)
	newContext := context.WithValue(ctx, LoggerKey, newLogger)

	return newContext
}
