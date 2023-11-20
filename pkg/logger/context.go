package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	loggerCtxKey = "logger"
)

func SetInCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func FromCtx(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return structuredLogger{
			zap: zap.NewNop(),
		}
	}

	return l
}

func NewCtx(ctx context.Context) context.Context {
	newCtx := context.Background()
	// Copy trace span to new ctx
	trace.ContextWithSpan(newCtx, trace.SpanFromContext(ctx))

	return SetInCtx(newCtx, FromCtx(ctx))
}
