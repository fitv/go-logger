package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// Logger is the logger wrapper.
type Logger struct {
	drivers    []Driver
	level      Level
	timeLayout string
}

// Driver is the logger interface.
type Driver interface {
	io.Writer
	Close() error
}

// Option is the option for logger.
type Option struct {
	// Path is the log file. eg: /var/log
	Path string

	// Filename  of the log file. eg: web
	Name string

	// Whether to create a new file every day.
	Daily bool

	// Days is the number of days to keep log files.
	Days int
}

// New creates a new Logger.
func New(drivers ...Driver) *Logger {
	if len(drivers) == 0 {
		drivers = append(drivers, NewStdLogger())
	}
	return &Logger{
		drivers:    drivers,
		level:      DebugLevel,
		timeLayout: "2006-01-02 15:04:05",
	}
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

// Fatal writes a message to the log using the FATAL level.
func (l *Logger) Fatal(args ...any) {
	l.Log(FatalLevel, args...)
}

// Log writes a message to the log using the given level.
func (l *Logger) Log(level Level, args ...any) {
	// Skip if the level is below the logger level.
	if level < l.level {
		return
	}

	args = append([]interface{}{
		time.Now().Format(l.timeLayout),
		fmt.Sprintf("%s:", strings.ToUpper(level.String())),
	}, args...)

	for i := range l.drivers {
		_, err := l.drivers[i].Write([]byte(fmt.Sprintln(args...)))
		// when write log to drive failed, print to stdout
		if err != nil {
			log.Println(fmt.Errorf("write log failed: %w", err))
		}
	}
}

// Close closes the logger.
func (l *Logger) Close() error {
	var err error
	for i := range l.drivers {
		if e := l.drivers[i].Close(); e != nil {
			err = e
		}
	}
	return err
}
