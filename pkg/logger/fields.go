package logger

import "go.uber.org/zap"

func WithString(key string, value string) Field {
	return field{
		zapField: zap.String(key, value),
	}
}

func WithInt(key string, value int) Field {
	return field{
		zapField: zap.Int(key, value),
	}
}
