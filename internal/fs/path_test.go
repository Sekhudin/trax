package fs

import (
	"path/filepath"
	"testing"
)

func TestOSPath_Success(t *testing.T) {
	p := NewPath()

	t.Run("join_path_elements", func(t *testing.T) {
		out := p.Join("a", "b", "c.txt")
		expected := filepath.Join("a", "b", "c.txt")
		if out != expected {
			t.Fatalf("wrong_join_result")
		}
	})

	t.Run("get_directory_path", func(t *testing.T) {
		path := filepath.Join("a", "b", "c.txt")
		if p.Dir(path) != filepath.Dir(path) {
			t.Fatal("wrong_dir_result")
		}
	})

	t.Run("get_base_name", func(t *testing.T) {
		if p.Base("a/b/c.txt") != "c.txt" {
			t.Fatal("wrong_base_result")
		}
	})

	t.Run("get_file_extension", func(t *testing.T) {
		if p.Ext("file.test.ts") != ".ts" {
			t.Fatal("wrong_ext_result")
		}
	})
}

func TestOSPath_Fallback(t *testing.T) {
	p := NewPath()

	t.Run("handle_no_extension", func(t *testing.T) {
		if p.Ext("Makefile") != "" {
			t.Fatal("expected_empty_ext")
		}
	})
}

func TestIsAllowedExt_Success(t *testing.T) {
	allowed := []string{".ts", ".tsx", ".js"}

	t.Run("allow_valid_extension", func(t *testing.T) {
		if !IsAllowedExt(".ts", allowed) {
			t.Fatal("should_be_allowed")
		}
	})
}

func TestIsAllowedExt_Fallback(t *testing.T) {
	allowed := []string{".ts", ".tsx", ".js"}

	t.Run("reject_invalid_extension", func(t *testing.T) {
		if IsAllowedExt(".go", allowed) {
			t.Fatal("should_be_rejected")
		}
	})

	t.Run("handle_empty_list", func(t *testing.T) {
		if IsAllowedExt(".ts", []string{}) {
			t.Fatal("should_be_false")
		}
	})

	t.Run("check_case_sensitivity", func(t *testing.T) {
		if IsAllowedExt(".TS", allowed) {
			t.Fatal("should_be_sensitive")
		}
	})
}
