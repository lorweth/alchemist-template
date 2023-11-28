package logger

type noopLogger struct{}

func NewNoop() Logger {
	return noopLogger{}
}

func (n noopLogger) Debugf(format string, args ...interface{}) {

}

func (n noopLogger) Infof(format string, args ...interface{}) {

}

func (n noopLogger) Warnf(format string, args ...interface{}) {

}

func (n noopLogger) Errorf(err error, format string, args ...interface{}) {

}

func (n noopLogger) With(fields ...Field) Logger {
	return n
}

func (n noopLogger) Flush() error {
	return nil
}

func (n noopLogger) clone() Logger {
	return n
}
