package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fitv/go-logger/internal/util"
)

var _ io.WriteCloser = (*FileWriter)(nil)

type FileWriter struct {
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

// NewFileWriter creates a new FileWriter.
func NewFileWriter(opt *Option) *FileWriter {
	ext := filepath.Ext(opt.Path)
	dir := filepath.Dir(opt.Path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	writer := &FileWriter{
		path:  opt.Path,
		daily: opt.Daily,
		days:  opt.Days,
		dir:   dir,
		ext:   ext,
		name:  filepath.Base(opt.Path[:len(opt.Path)-len(ext)]),
	}
	if writer.daily {
		writer.date = util.Today()
		writer.regexDateName = regexp.MustCompile(fmt.Sprintf(
			`%s-(\d{4}-\d{2}-\d{2})%s`,
			regexp.QuoteMeta(writer.name),
			regexp.QuoteMeta(writer.ext),
		))
	}
	return writer
}

// Write writes the log message.
func (l *FileWriter) Write(p []byte) (n int, err error) {
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
func (l *FileWriter) Close() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.close()
}

// close closes the log file.
func (l *FileWriter) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// openFile opens the log file.
func (l *FileWriter) openFile() error {
	file, err := os.OpenFile(l.filePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	l.file = file
	return nil
}

// filePath returns the log file full path.
func (l *FileWriter) filePath() string {
	if l.daily {
		return fmt.Sprintf("%s/%s-%s", l.dir, l.name, l.date+l.ext)
	}
	return l.path
}

// cleanOutdatedFiles deletes outdated log files.
func (l *FileWriter) cleanOutdatedFiles() error {
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
