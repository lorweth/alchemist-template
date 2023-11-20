package logger

// Logger represents a structuredLogger used for application logging.
type Logger interface {
	// Debugf logs a message at the debug level with optional message arguments.
	Debugf(format string, args ...interface{})

	// Infof logs a message at the info level with optional message arguments.
	Infof(format string, args ...interface{})

	// Warnf logs a message at the warning level with optional message arguments.
	Warnf(format string, args ...interface{})

	// Errorf logs an error with a description at the error level.
	Errorf(err error, format string, args ...interface{})

	// With creates a new structuredLogger instance with additional fields.
	With(fields ...Field) Logger

	// Flush flushing any buffered log entries
	Flush() error

	// clone creates a copy of Logger
	clone() Logger
}

// Field represents a field that can be added to a structuredLogger to provide additional context.
type Field interface {
	// Equals checks if this field is equal to another field, typically used for deduplication.
	Equals(other Field) bool
}
