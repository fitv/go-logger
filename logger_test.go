package logger_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fitv/go-logger"
	"github.com/fitv/go-logger/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFileLogger(t *testing.T) {
	assert := assert.New(t)

	filedLogger := logger.NewFileLogger(&logger.Option{
		Path:  "/tmp",
		Name:  "test",
		Daily: false,
	})

	log := logger.New(filedLogger)
	defer log.Close()
	log.SetLevel(logger.DebugLevel)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	log.Fatal("fatal")

	logPath := "/tmp/test.log"
	defer os.Remove(logPath)

	bytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}

	assert.Contains(string(bytes), "DEBUG: debug")
	assert.Contains(string(bytes), "INFO: info")
	assert.Contains(string(bytes), "WARN: warn")
	assert.Contains(string(bytes), "ERROR: error")
	assert.Contains(string(bytes), "FATAL: fatal")
}

func TestFileLoggerDaily(t *testing.T) {
	assert := assert.New(t)

	fileLogger := logger.NewFileLogger(&logger.Option{
		Path:  "/tmp",
		Name:  "test",
		Daily: true,
	})

	log := logger.New(fileLogger)
	defer log.Close()
	log.SetLevel(logger.WarnLevel)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	log.Fatal("fatal")

	date := util.Today()
	logPath := fmt.Sprintf("/tmp/test-%s.log", date)
	defer os.Remove(logPath)

	bytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}

	assert.NotContains(string(bytes), "DEBUG: debug")
	assert.NotContains(string(bytes), "INFO: info")
	assert.Contains(string(bytes), "WARN: warn")
	assert.Contains(string(bytes), "ERROR: error")
	assert.Contains(string(bytes), "FATAL: fatal")
}
