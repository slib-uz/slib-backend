package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const loggerKey ctxKey = iota

func WithLogger(ctx context.Context, lg *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, lg)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return nil
	}
	if lg, ok := ctx.Value(loggerKey).(*zap.Logger); ok && lg != nil {
		return lg
	}
	return nil
}
