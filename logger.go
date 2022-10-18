package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var _ io.Writer = (*Logger)(nil)

type Option struct {
	// Path of the log file. eg: /logs/app.log
	Path string

	// Whether to create a new file every day.
	Daily bool

	// Days is the number of days to keep log files.
	Days int
}

type Logger struct {
	out        io.Writer
	level      Level
	timeLayout string
}

// New creates a new Logger.
func New() *Logger {
	return &Logger{
		out:        os.Stdout,
		level:      InfoLevel,
		timeLayout: "2006-01-02 15:04:05",
	}
}

// SetOut sets the output writer.
func (l *Logger) SetOutput(out io.Writer) {
	l.out = out
}

// SetLevel sets the logger level.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetTimeLayout sets the time layout.
func (l *Logger) SetTimeLayout(layout string) {
	l.timeLayout = layout
}

// Debug writes a message to the log using the DEBUG level.
func (l *Logger) Debug(args ...any) {
	l.Log(DebugLevel, args...)
}

// Info writes a message to the log using the INFO level.
func (l *Logger) Info(args ...any) {
	l.Log(InfoLevel, args...)
}

// Warn writes a message to the log using the WARN level.
func (l *Logger) Warn(args ...any) {
	l.Log(WarnLevel, args...)
}

// Error writes a message to the log using the ERROR level.
func (l *Logger) Error(args ...any) {
	l.Log(ErrorLevel, args...)
}

// Panic writes a message to the log using the PANIC level. The process will panic after writing the message.
func (l *Logger) Panic(args ...any) {
	l.Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Fatal writes a message to the log using the FATAL level. The process will exit with status set to 1.
func (l *Logger) Fatal(args ...any) {
	l.Log(FatalLevel, args...)
	os.Exit(1)
}

func (l *Logger) Write(p []byte) (n int, err error) {
	return l.out.Write(l.format(InfoLevel, string(p)))
}

// Log writes a message to the log using the given level.
func (l *Logger) Log(level Level, args ...any) {
	if level < l.level {
		return
	}

	_, err := l.out.Write(l.format(level, fmt.Sprint(args...)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write log failed: %v\n", err)
	}
}

// format returns a formatted log message.
func (l *Logger) format(level Level, msg string) []byte {
	b := &bytes.Buffer{}

	b.WriteString(time.Now().Format(l.timeLayout))
	b.WriteByte(' ')
	b.WriteString(strings.ToUpper(level.String()))
	b.WriteByte(':')
	b.WriteByte(' ')
	b.WriteString(msg)
	b.WriteByte('\n')

	return b.Bytes()
}
