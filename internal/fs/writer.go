package fs

import (
	"os"
	"path/filepath"
)

type Writer interface {
	Write(path string, data []byte) error
}

type OSWriter struct{}

func NewWriter() *OSWriter {
	return &OSWriter{}
}

func (w *OSWriter) Write(path string, data []byte) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
