package fs

import (
	"fmt"
	"os"
	"path/filepath"

	appErr "trax/internal/errors"
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
		return appErr.NewIOError("path", fmt.Sprintf("failed to create directory %q", path), err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return appErr.NewIOError("path", fmt.Sprintf("failed to write file %q", path), err)
	}

	return nil
}
