package fs

import "os"

type Exister interface {
	Exists(path string) bool
}

type OSExister struct{}

func NewExister() *OSExister {
	return &OSExister{}
}

func (e *OSExister) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
