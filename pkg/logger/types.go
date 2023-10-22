package logger

// Logger represents a logger used for application logging.
type Logger interface {
	// Debug logs a message at the debug level with optional message arguments.
	Debug(msg string, args ...interface{})

	// Info logs a message at the info level with optional message arguments.
	Info(msg string, args ...interface{})

	// Warn logs a message at the warning level with optional message arguments.
	Warn(msg string, args ...interface{})

	// Error logs an error with a description at the error level.
	Error(err error, desc string)

	// With creates a new logger instance with additional fields.
	With(fields ...Field) Logger
}

// Field represents a field that can be added to a logger to provide additional context.
type Field interface {
	// Equals checks if this field is equal to another field, typically used for deduplication.
	Equals(other Field) bool
}
