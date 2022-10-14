package logger

import "os"

var _ Driver = (*StdLogger)(nil)

// FileLogger is a file logger struct.
type StdLogger struct{}

// NewFileLogger creates a new FileLogger.
func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

// WithFields adds fields to the logger.
func (l *StdLogger) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// Close closes the logger.
func (l *StdLogger) Close() error {
	return nil
}
