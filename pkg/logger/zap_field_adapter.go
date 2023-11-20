package logger

import (
	"go.uber.org/zap"
)

type field struct {
	zapField zap.Field
}

func (f field) Equals(other Field) bool {
	otherField, ok := other.(field)
	if !ok {
		return false
	}

	return f.zapField.Equals(otherField.zapField)
}

func (f field) Unwrap() zap.Field {
	return f.zapField
}
