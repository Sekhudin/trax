package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestOSWriter_Write(t *testing.T) {
	w := NewOSWriter()

	t.Run("should create directory and write file successfully", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "nested", "file.txt")

		err := w.Write(fp, []byte("hello"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, err := os.ReadFile(fp)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}

		if string(b) != "hello" {
			t.Fatalf("unexpected file content: %q", string(b))
		}
	})

	t.Run("should overwrite existing file", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")

		if err := os.WriteFile(fp, []byte("old"), 0o644); err != nil {
			t.Fatalf("failed to seed file: %v", err)
		}

		err := w.Write(fp, []byte("new"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, _ := os.ReadFile(fp)
		if string(b) != "new" {
			t.Fatalf("file not overwritten")
		}
	})

	t.Run("should return IO error when directory cannot be created", func(t *testing.T) {
		dir := t.TempDir()
		block := filepath.Join(dir, "blocked")

		if err := os.WriteFile(block, []byte("x"), 0o644); err != nil {
			t.Fatalf("failed to create blocking file: %v", err)
		}

		fp := filepath.Join(block, "file.txt")

		err := w.Write(fp, []byte("data"))
		if err == nil {
			t.Fatalf("expected error")
		}

		var ce *appErr.CoreError
		if !errors.As(err, &ce) {
			t.Fatalf("expected CoreError")
		}

		if ce.Code != appErr.ErrIO {
			t.Fatalf("expected ErrIO, got %s", ce.Code)
		}
	})

	t.Run("should return IO error when file cannot be written", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "readonly")

		if err := os.Mkdir(sub, 0o555); err != nil {
			t.Fatalf("failed to create readonly dir: %v", err)
		}
		defer os.Chmod(sub, 0o755)

		fp := filepath.Join(sub, "file.txt")

		err := w.Write(fp, []byte("data"))
		if err == nil {
			t.Fatalf("expected write error")
		}

		var ce *appErr.CoreError
		if !errors.As(err, &ce) {
			t.Fatalf("expected CoreError")
		}

		if ce.Code != appErr.ErrIO {
			t.Fatalf("expected ErrIO, got %s", ce.Code)
		}
	})
}
