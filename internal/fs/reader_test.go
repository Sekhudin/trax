package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSReader_Success(t *testing.T) {
	r := NewReader()

	t.Run("read_file_content", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")
		content := []byte("hello world")

		os.WriteFile(fp, content, 0o644)

		b, err := r.Read(fp)
		if err != nil {
			t.Fatalf("unexpected_read_error: %v", err)
		}

		if string(b) != "hello world" {
			t.Fatalf("content_mismatch: %q", string(b))
		}
	})
}

func TestOSReader_Error(t *testing.T) {
	r := NewReader()

	t.Run("file_not_found", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "missing.txt")

		if _, err := r.Read(fp); err == nil {
			t.Fatal("expected_missing_error")
		}
	})

	t.Run("is_a_directory", func(t *testing.T) {
		dir := t.TempDir()
		if _, err := r.Read(dir); err == nil {
			t.Fatal("expected_directory_error")
		}
	})

	t.Run("permission_denied", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "no_access.txt")

		os.WriteFile(fp, []byte("secret"), 0o000)
		defer os.Chmod(fp, 0o644)

		if _, err := r.Read(fp); err == nil {
			t.Fatal("expected_permission_error")
		}
	})
}
