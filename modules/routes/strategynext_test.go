package routes

import (
	"io/fs"
	"path/filepath"
	"testing"
)

type fakeDir struct {
	name string
	dir  bool
}

func (f fakeDir) Name() string               { return f.name }
func (f fakeDir) IsDir() bool                { return f.dir }
func (f fakeDir) Type() fs.FileMode          { return 0 }
func (f fakeDir) Info() (fs.FileInfo, error) { return nil, nil }

func TestNextRule_AllBranches(t *testing.T) {
	r := newNextRule()

	t.Run("isNonRouteDir", func(t *testing.T) {
		if !r.isNonRouteDir(fakeDir{"_private", true}) {
			t.Fatal("expected non route dir")
		}
		if r.isNonRouteDir(fakeDir{"file", false}) {
			t.Fatal("file should not be dir")
		}
	})

	t.Run("isSlotDir", func(t *testing.T) {
		if !r.isSlotDir(fakeDir{"@slot", true}) {
			t.Fatal("expected slot dir")
		}
		if r.isSlotDir(fakeDir{"users", true}) {
			t.Fatal("should be false for normal dir")
		}
		if r.isSlotDir(fakeDir{"@file", false}) {
			t.Fatal("file should not be treated as slot dir")
		}
	})

	t.Run("shouldSkip", func(t *testing.T) {
		cases := []string{
			"api",
			"(.)test",
			"(..)test",
			"(..)(..)test",
			"(...)test",
		}

		for _, c := range cases {
			if r.shouldSkip(fakeDir{c, true}) != filepath.SkipDir {
				t.Fatalf("should skip %s", c)
			}
		}

		if r.shouldSkip(fakeDir{"users", true}) != nil {
			t.Fatal("should not skip")
		}
	})

	t.Run("segmentKind", func(t *testing.T) {
		tests := map[string]string{
			"(group)":     "Group",
			"[[...slug]]": "OSlug",
			"[...slug]":   "Slug",
			"[id]":        "Params",
			"users":       "Static",
		}

		for in, want := range tests {
			if r.segmentKind(in) != want {
				t.Fatalf("%s failed", in)
			}
		}
	})

	t.Run("normalizeSegment success", func(t *testing.T) {
		tests := map[string]string{
			"":            "",
			"(group)":     "",
			"[id]":        ":id",
			"[...slug]":   "*",
			"[[...slug]]": "*",
			"users":       "users",
		}

		for in, want := range tests {
			got, err := r.normalizeSegment(in)
			if err != nil || got != want {
				t.Fatalf("%s failed", in)
			}
		}
	})

	t.Run("normalizeSegment error", func(t *testing.T) {
		if _, err := r.normalizeSegment("[id"); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("getAffix valid", func(t *testing.T) {
		a := r.getAffix("Params")
		if a.pre != "[" {
			t.Fatal("wrong affix")
		}
	})

	t.Run("getAffix panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		r.getAffix("NotExist")
	})
}

func TestNextAppAndNextPage_ShouldSkip(t *testing.T) {
	app := newNextApp()
	page := newNextPage()

	t.Run("nextapp non route", func(t *testing.T) {
		if app.shouldSkip("", fakeDir{"_hidden", true}) != filepath.SkipDir {
			t.Fatal("nextapp non route failed")
		}
	})

	t.Run("nextapp slot", func(t *testing.T) {
		if app.shouldSkip("", fakeDir{"@slot", true}) != filepath.SkipDir {
			t.Fatal("nextapp slot failed")
		}
	})

	t.Run("nextapp rule skip", func(t *testing.T) {
		if app.shouldSkip("", fakeDir{"api", true}) != filepath.SkipDir {
			t.Fatal("nextapp rule skip failed")
		}
	})

	t.Run("nextpage rule skip", func(t *testing.T) {
		if page.shouldSkip("", fakeDir{"api", true}) != filepath.SkipDir {
			t.Fatal("nextpage skip failed")
		}
	})

	t.Run("nextapp normalizeSegment passthrough", func(t *testing.T) {
		s, err := app.normalizeSegment("[id]")
		if err != nil || s != ":id" {
			t.Fatal("nextapp normalizeSegment failed")
		}
	})

	t.Run("nextpage normalizeSegment passthrough", func(t *testing.T) {
		s, err := page.normalizeSegment("[...slug]")
		if err != nil || s != "*" {
			t.Fatal("nextpage normalizeSegment failed")
		}
	})
}
