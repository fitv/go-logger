package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fitv/go-logger/internal/util"
)

var _ Driver = (*FileLogger)(nil)

// FileLogger is a file logger struct.
type FileLogger struct {
	mux           sync.Mutex
	file          *os.File
	name          string
	ext           string
	dir           string
	path          string
	date          string
	daily         bool
	days          int
	regexDateName *regexp.Regexp
}

// NewFileLogger creates a new FileLogger.
func NewFileLogger(opt *Option) *FileLogger {
	ext := filepath.Ext(opt.Path)

	logger := &FileLogger{
		ext:   ext,
		path:  opt.Path,
		daily: opt.Daily,
		days:  opt.Days,
		dir:   filepath.Dir(opt.Path),
		name:  filepath.Base(opt.Path[:len(opt.Path)-len(ext)]),
	}
	if logger.daily {
		logger.date = util.Today()
		logger.regexDateName = regexp.MustCompile(fmt.Sprintf(
			`%s-(\d{4}-\d{2}-\d{2})%s`,
			regexp.QuoteMeta(logger.name),
			regexp.QuoteMeta(logger.ext),
		))
	}
	return logger
}

// WithFields adds fields to the logger.
func (l *FileLogger) Write(p []byte) (n int, err error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.file == nil {
		if err := l.openFile(); err != nil {
			return 0, err
		}
	}
	if l.daily && util.Today() != l.date {
		if err := l.close(); err != nil {
			return 0, err
		}
		l.date = util.Today()
		if err := l.openFile(); err != nil {
			l.date = ""
			return 0, err
		}
		if err := l.cleanOutdatedFiles(); err != nil {
			return 0, err
		}
	}

	return l.file.Write(p)
}

// Close closes the logger.
func (l *FileLogger) Close() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.close()
}

// close closes the log file.
func (l *FileLogger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// openFile opens the log file.
func (l *FileLogger) openFile() error {
	file, err := os.OpenFile(l.filePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	l.file = file
	return nil
}

// filePath returns the log file full path.
func (l *FileLogger) filePath() string {
	if l.daily {
		return fmt.Sprintf("%s/%s-%s", l.dir, l.name, l.date+l.ext)
	}
	return l.path
}

// cleanOutdatedFiles deletes outdated log files.
func (l *FileLogger) cleanOutdatedFiles() error {
	if !(l.daily && l.days > 0) {
		return nil
	}

	dirEntries, err := os.ReadDir(l.dir)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		matches := l.regexDateName.FindStringSubmatch(dirEntry.Name())
		if !(len(matches) > 1 && util.IsValidDate(matches[1]) && util.DiffDays(matches[1]) > l.days) {
			continue
		}
		if err := os.Remove(filepath.Join(l.dir, dirEntry.Name())); err != nil {
			return err
		}
	}
	return nil
}
