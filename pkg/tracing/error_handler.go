package tracing

import (
	"github.com/virsavik/alchemist-template/pkg/logger"
)

type ErrorHandler struct {
	logger logger.Logger
}

func NewErrorHandler(log logger.Logger) ErrorHandler {
	return ErrorHandler{
		logger: log,
	}
}

func (e ErrorHandler) Handle(err error) {
	e.logger.Errorf(err, "otel err:")
}
