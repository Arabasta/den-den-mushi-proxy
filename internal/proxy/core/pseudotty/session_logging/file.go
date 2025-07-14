package session_logging

import (
	"io"
	"os"
	"path/filepath"
)

type FileSessionLogger struct {
	writer io.WriteCloser
}

func NewFileSessionLogger(path string) (*FileSessionLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &FileSessionLogger{writer: f}, nil
}

func (f *FileSessionLogger) WriteLine(line string) error {
	_, err := f.writer.Write([]byte(line + "\n"))
	return err
}

func (f *FileSessionLogger) Close() error {
	return f.writer.Close()
}
