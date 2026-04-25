package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSExister_Exists(t *testing.T) {
	ex := NewExister()

	t.Run("should return true when file exists", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")

		if err := os.WriteFile(fp, []byte("hello"), 0o644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		if !ex.Exists(fp) {
			t.Fatalf("expected file to exist")
		}
	})

	t.Run("should return true when path is a directory", func(t *testing.T) {
		dir := t.TempDir()

		if !ex.Exists(dir) {
			t.Fatalf("expected directory to exist")
		}
	})

	t.Run("should return false when path does not exist", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "notfound.txt")

		if ex.Exists(fp) {
			t.Fatalf("expected path to not exist")
		}
	})

	t.Run("should return false when no permission to stat path", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "restricted")

		if err := os.Mkdir(sub, 0o000); err != nil {
			t.Fatalf("failed to create restricted dir: %v", err)
		}
		defer os.Chmod(sub, 0o755)

		fp := filepath.Join(sub, "file.txt")

		if ex.Exists(fp) {
			t.Fatalf("expected false due to permission issue")
		}
	})
}
