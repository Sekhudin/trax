package fs

import "os"

type Reader interface {
	Read(path string) ([]byte, error)
}

type OSReader struct{}

func NewReader() *OSReader {
	return &OSReader{}
}

func (r *OSReader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
