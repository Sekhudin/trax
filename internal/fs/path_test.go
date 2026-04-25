package fs

import (
	"path/filepath"
	"testing"
)

func TestOSPath_Methods(t *testing.T) {
	p := NewPath()

	t.Run("Join should combine paths correctly", func(t *testing.T) {
		out := p.Join("a", "b", "c.txt")
		expected := filepath.Join("a", "b", "c.txt")

		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("Dir should return directory of path", func(t *testing.T) {
		path := filepath.Join("a", "b", "c.txt")
		out := p.Dir(path)
		expected := filepath.Dir(path)

		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("Base should return last element", func(t *testing.T) {
		path := filepath.Join("a", "b", "c.txt")
		out := p.Base(path)
		expected := "c.txt"

		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("Ext should return file extension", func(t *testing.T) {
		path := "file.test.ts"
		out := p.Ext(path)

		if out != ".ts" {
			t.Fatalf("expected .ts, got %q", out)
		}
	})

	t.Run("Ext should return empty string when no extension", func(t *testing.T) {
		out := p.Ext("Makefile")

		if out != "" {
			t.Fatalf("expected empty extension, got %q", out)
		}
	})
}

func TestIsAllowedExt(t *testing.T) {
	allowed := []string{".ts", ".tsx", ".js"}

	t.Run("should return true when extension is allowed", func(t *testing.T) {
		if !IsAllowedExt(".ts", allowed) {
			t.Fatalf("expected extension to be allowed")
		}
	})

	t.Run("should return false when extension is not allowed", func(t *testing.T) {
		if IsAllowedExt(".go", allowed) {
			t.Fatalf("expected extension to not be allowed")
		}
	})

	t.Run("should return false when allowed list is empty", func(t *testing.T) {
		if IsAllowedExt(".ts", []string{}) {
			t.Fatalf("expected false when allowed list is empty")
		}
	})

	t.Run("should be case sensitive", func(t *testing.T) {
		if IsAllowedExt(".TS", allowed) {
			t.Fatalf("expected case sensitive check")
		}
	})
}
