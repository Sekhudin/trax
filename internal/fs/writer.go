package fs

import (
	"os"
	"path/filepath"
)

type FileWriter interface {
	Write(path string, data []byte) error
}

type OSWriter struct{}

func NewOSWriter() FileWriter {
	return &OSWriter{}
}

func (w *OSWriter) Write(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
