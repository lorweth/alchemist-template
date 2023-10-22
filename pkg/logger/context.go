package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	loggerCtxKey = "logger"
)

func NewCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func FromCtx(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return logger{
			zap: zap.NewNop(),
		}
	}

	return l
}
