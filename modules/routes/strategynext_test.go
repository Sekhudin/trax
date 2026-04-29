package routes

import (
	"io/fs"
	"path/filepath"
	"testing"
)

type mockDir struct {
	name string
	dir  bool
}

func (f mockDir) Name() string               { return f.name }
func (f mockDir) IsDir() bool                { return f.dir }
func (f mockDir) Type() fs.FileMode          { return 0 }
func (f mockDir) Info() (fs.FileInfo, error) { return nil, nil }

func TestNextRule_Success(t *testing.T) {
	r := newNextRule()

	t.Run("special_dir_detection", func(t *testing.T) {
		if !r.isNonRouteDir(mockDir{"_private", true}) {
			t.Fatal("should_detect_non_route")
		}
		if !r.isSlotDir(mockDir{"@slot", true}) {
			t.Fatal("should_detect_slot")
		}
	})

	t.Run("should_skip_folders", func(t *testing.T) {
		cases := []string{"api", "(.)test", "(..)test", "(..)(..)test", "(...)test"}
		for _, c := range cases {
			if r.shouldSkip(mockDir{c, true}) != filepath.SkipDir {
				t.Fatalf("should_skip_%s", c)
			}
		}
		if r.shouldSkip(mockDir{"users", true}) != nil {
			t.Fatal("should_not_skip_normal_dir")
		}
	})

	t.Run("normalize_segments_kinds", func(t *testing.T) {
		tests := map[string]string{
			"":            "",
			"(group)":     "",
			"[id]":        ":id",
			"[...slug]":   "*",
			"[[...slug]]": "*",
			"static":      "static",
		}

		for in, want := range tests {
			got := r.normalizeSegment(in)
			if got != want {
				t.Fatalf("in: %s, want: %s, got: %s", in, want, got)
			}
		}
	})
}

func TestNextRule_Error(t *testing.T) {
	r := newNextRule()

	t.Run("file_is_not_special_dir", func(t *testing.T) {
		if r.isNonRouteDir(mockDir{"_file", false}) || r.isSlotDir(mockDir{"@file", false}) {
			t.Fatal("should_ignore_if_not_dir")
		}
	})
}

func TestNextApp_Success(t *testing.T) {
	app := newNextApp()

	t.Run("normalize_segment", func(t *testing.T) {
		res := app.normalizeSegment("[id]")
		if res != ":id" {
			t.Fatal("should_normalize_to_params")
		}
	})

	t.Run("should_skip_dirs", func(t *testing.T) {
		if app.shouldSkip("", mockDir{"_any", true}) != filepath.SkipDir {
			t.Fatal("should_skip_non_route_dir")
		}

		if app.shouldSkip("", mockDir{"@any", true}) != filepath.SkipDir {
			t.Fatal("should_skip_slot_dir")
		}

		if app.shouldSkip("", mockDir{"api", true}) != filepath.SkipDir {
			t.Fatal("should_skip_forbidden_folder")
		}
	})
}

func TestNextPage_Success(t *testing.T) {
	page := newNextPage()

	t.Run("normalize_segment", func(t *testing.T) {
		res := page.normalizeSegment("[...slug]")
		if res != "*" {
			t.Fatal("should_normalize_to_wildcard")
		}
	})

	t.Run("should_skip_dirs", func(t *testing.T) {
		if page.shouldSkip("", mockDir{"api", true}) != filepath.SkipDir {
			t.Fatal("should_skip_forbidden_folder")
		}
	})
}
