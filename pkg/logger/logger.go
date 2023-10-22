package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type LogConfig struct {
	Environment string
}

type logger struct {
	zap *zap.Logger
}

func New(cfg LogConfig) (Logger, error) {
	var zapLogger *zap.Logger
	var err error

	switch cfg.Environment {
	case "PRODUCTION":
		zapLogger, err = newZapConfig(withProductionConfig()).Build()
	default:
		zapLogger, err = newZapConfig().Build()
	}

	if err != nil {
		return &logger{
			zap: zap.NewNop(),
		}, nil
	}

	return &logger{
		zap: zapLogger,
	}, nil

}

func (l logger) Debug(msg string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(msg, args...))
}

func (l logger) Error(err error, desc string) {
	l.zap.Error(fmt.Sprintf("%s %v", desc, err))
}

func (l logger) Info(msg string, args ...interface{}) {
	l.zap.Info(fmt.Sprintf(msg, args...))
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.zap.Warn(fmt.Sprintf(msg, args...))
}

func (l logger) With(fields ...Field) Logger {
	zapFields := make([]zap.Field, len(fields))
	for idx, f := range fields {
		zapFields[idx] = f.(field).Unwrap()
	}

	return &logger{
		zap: l.zap.With(zapFields...),
	}
}
