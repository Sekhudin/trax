package fs

import (
	"path/filepath"
	"slices"
)

type Path interface {
	Join(elem ...string) string
	Dir(path string) string
	Base(path string) string
	Ext(path string) string
}

type OSPath struct{}

func NewPath() *OSPath {
	return &OSPath{}
}

func (p *OSPath) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (p *OSPath) Dir(path string) string {
	return filepath.Dir(path)
}

func (p *OSPath) Base(path string) string {
	return filepath.Base(path)
}

func (p *OSPath) Ext(path string) string {
	return filepath.Ext(path)
}

func IsAllowedExt(ext string, allowed []string) bool {
	return slices.Contains(allowed, ext)
}
