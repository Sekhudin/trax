package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSExister_Success(t *testing.T) {
	ex := NewExister()

	t.Run("detect_existing_file", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "file.txt")
		os.WriteFile(fp, []byte("hello"), 0o644)

		if !ex.Exists(fp) {
			t.Fatal("failed_detect_file")
		}
	})

	t.Run("detect_existing_directory", func(t *testing.T) {
		dir := t.TempDir()
		if !ex.Exists(dir) {
			t.Fatal("failed_detect_dir")
		}
	})
}

func TestOSExister_Fallback(t *testing.T) {
	ex := NewExister()

	t.Run("handle_missing_path", func(t *testing.T) {
		dir := t.TempDir()
		fp := filepath.Join(dir, "notfound.txt")
		if ex.Exists(fp) {
			t.Fatal("should_return_false")
		}
	})

	t.Run("handle_permission_denied", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "restricted")
		os.Mkdir(sub, 0o000)
		defer os.Chmod(sub, 0o755)

		fp := filepath.Join(sub, "file.txt")
		if ex.Exists(fp) {
			t.Fatal("should_fail_permission")
		}
	})
}
