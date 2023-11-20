package logger

import (
	"go.uber.org/zap"
)

func String(key string, value string) Field {
	return field{
		zapField: zap.String(key, value),
	}
}

func Int(key string, value int) Field {
	return field{
		zapField: zap.Int(key, value),
	}
}
