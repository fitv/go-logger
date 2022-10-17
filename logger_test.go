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

	path := "/tmp/test.log"
	defer os.Remove(path)

	filedLogger := logger.NewFileLogger(&logger.Option{
		Path:  path,
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

	bytes, err := os.ReadFile(path)
	assert.NoError(err)

	assert.Contains(string(bytes), "DEBUG: debug")
	assert.Contains(string(bytes), "INFO: info")
	assert.Contains(string(bytes), "WARN: warn")
	assert.Contains(string(bytes), "ERROR: error")
	assert.Contains(string(bytes), "FATAL: fatal")
}

func TestFileLoggerDaily(t *testing.T) {
	assert := assert.New(t)

	fileLogger := logger.NewFileLogger(&logger.Option{
		Path:  "/tmp/test.log",
		Daily: true,
		Days:  3,
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
	assert.NoError(err)

	assert.NotContains(string(bytes), "DEBUG: debug")
	assert.NotContains(string(bytes), "INFO: info")
	assert.Contains(string(bytes), "WARN: warn")
	assert.Contains(string(bytes), "ERROR: error")
	assert.Contains(string(bytes), "FATAL: fatal")
}
