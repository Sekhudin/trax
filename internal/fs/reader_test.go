package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSReader_Read(t *testing.T) {
	r := NewReader()

	t.Run("should read file successfully", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")
		content := []byte("hello world")

		if err := os.WriteFile(fp, content, 0o644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		b, err := r.Read(fp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(b) != "hello world" {
			t.Fatalf("expected content 'hello world', got %q", string(b))
		}
	})

	t.Run("should return error when file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "missing.txt")

		_, err := r.Read(fp)
		if err == nil {
			t.Fatalf("expected error for missing file")
		}
	})

	t.Run("should return error when path is a directory", func(t *testing.T) {
		dir := t.TempDir()

		_, err := r.Read(dir)
		if err == nil {
			t.Fatalf("expected error when reading a directory")
		}
	})

	t.Run("should return error when permission denied", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "restricted.txt")

		if err := os.WriteFile(fp, []byte("secret"), 0o000); err != nil {
			t.Fatalf("failed to create restricted file: %v", err)
		}
		defer os.Chmod(fp, 0o644)

		_, err := r.Read(fp)
		if err == nil {
			t.Fatalf("expected permission error")
		}
	})
}
