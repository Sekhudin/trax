package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestOSWriter_Success(t *testing.T) {
	w := NewOSWriter()

	t.Run("write_nested_file", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "nested", "file.txt")

		if err := w.Write(fp, []byte("hello")); err != nil {
			t.Fatalf("unexpected_write_error: %v", err)
		}

		b, _ := os.ReadFile(fp)
		if string(b) != "hello" {
			t.Fatal("content_mismatch")
		}
	})

	t.Run("overwrite_existing_file", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")
		os.WriteFile(fp, []byte("old"), 0o644)

		if err := w.Write(fp, []byte("new")); err != nil {
			t.Fatal(err)
		}

		b, _ := os.ReadFile(fp)
		if string(b) != "new" {
			t.Fatal("file_not_overwritten")
		}
	})
}

func TestOSWriter_Error(t *testing.T) {
	w := NewOSWriter()

	t.Run("fail_create_directory", func(t *testing.T) {
		dir := t.TempDir()
		block := filepath.Join(dir, "blocked")
		os.WriteFile(block, []byte("x"), 0o644)

		fp := filepath.Join(block, "file.txt")
		err := w.Write(fp, []byte("data"))

		var ce *appErr.CoreError
		if !errors.As(err, &ce) || ce.Code != appErr.ErrIO {
			t.Fatal("expected_io_error_for_dir")
		}
	})

	t.Run("fail_write_file", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "readonly")
		os.Mkdir(sub, 0o555)
		defer os.Chmod(sub, 0o755)

		fp := filepath.Join(sub, "file.txt")
		err := w.Write(fp, []byte("data"))

		var ce *appErr.CoreError
		if !errors.As(err, &ce) || ce.Code != appErr.ErrIO {
			t.Fatal("expected_io_error_for_file")
		}
	})
}
