package logger

import (
	"errors"
	"fmt"
	"syscall"

	"go.uber.org/zap"
)

type Config struct {
	Environment string
}

type structuredLogger struct {
	zap *zap.Logger
}

func New(cfg Config) (Logger, error) {
	var zapLogger *zap.Logger
	var err error

	switch cfg.Environment {
	case "PRODUCTION":
		zapLogger, err = newZapConfig(withProductionConfig()).Build()
	default:
		zapLogger, err = newZapConfig().Build()
	}

	if err != nil {
		return &structuredLogger{
			zap: zap.NewNop(),
		}, nil
	}

	return &structuredLogger{
		zap: zapLogger,
	}, nil
}

func (l structuredLogger) Debugf(format string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, args...))
}

func (l structuredLogger) Errorf(err error, format string, args ...interface{}) {
	// Append otel exception keys
	zapLogger := l.zap.With(
		zap.String("exception.type", fmt.Sprintf("%T", err)),
		zap.String("exception.message", err.Error()),
	)

	zapLogger.Error(fmt.Sprintf(format+" %+v", append(args, err)...))
}

func (l structuredLogger) Infof(format string, args ...interface{}) {
	l.zap.Info(fmt.Sprintf(format, args...))
}

func (l structuredLogger) Warnf(format string, args ...interface{}) {
	l.zap.Warn(fmt.Sprintf(format, args...))
}

func (l structuredLogger) With(fields ...Field) Logger {
	zapFields := make([]zap.Field, len(fields))
	for idx, f := range fields {
		zapFields[idx] = f.(field).Unwrap()
	}

	return &structuredLogger{
		zap: l.zap.With(zapFields...),
	}
}

func (l structuredLogger) Flush() error {
	if err := l.zap.Sync(); err != nil {
		// Ignore this stderr https://github.com/uber-go/zap/issues/328
		if !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) {
			return err
		}
	}

	return nil
}

func (l structuredLogger) clone() Logger {
	zapCloned := *l.zap // Copy value of l.zap

	return structuredLogger{
		zap: &zapCloned,
	}
}
