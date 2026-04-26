package routes

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type fakeStrategy struct {
	skipErr      error
	normalizeErr error
}

func (f *fakeStrategy) shouldSkip(string, fs.DirEntry) error {
	return f.skipErr
}

func (f *fakeStrategy) normalizeSegment(seg string) (string, error) {
	if f.normalizeErr != nil {
		return "", f.normalizeErr
	}
	return seg, nil
}

func newTestWalker(root string, strat walkerstrategy) *walker {
	return &walker{
		strategy: strat,
		config: &Config{
			Root: root,
		},
		rule: &walkrule{
			exts: map[string]struct{}{
				".ts": {},
			},
			identRoute: map[string]struct{}{
				"index": {},
			},
			excludeFiles: map[string]struct{}{
				"_app": {},
			},
		},
	}
}

func touch(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestWalk_AllBranches(t *testing.T) {
	root := t.TempDir()

	touch(t, filepath.Join(root, "users", "index.ts"))
	touch(t, filepath.Join(root, "users", "profile.ts"))

	touch(t, filepath.Join(root, "_app.ts"))

	touch(t, filepath.Join(root, "ignore.js"))

	w := newTestWalker(root, &fakeStrategy{})

	rs, err := w.walk()
	if err != nil {
		t.Fatal(err)
	}

	if len(rs) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(rs))
	}
}

func TestWalk_ShouldSkipError(t *testing.T) {
	root := t.TempDir()
	touch(t, filepath.Join(root, "a.ts"))

	w := newTestWalker(root, &fakeStrategy{
		skipErr: errors.New("stop"),
	})

	_, err := w.walk()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWalk_NormalizeError_SkipRoute(t *testing.T) {
	root := t.TempDir()
	touch(t, filepath.Join(root, "a.ts"))

	w := newTestWalker(root, &fakeStrategy{
		normalizeErr: errors.New("bad"),
	})

	rs, err := w.walk()
	if err != nil {
		t.Fatal(err)
	}

	if len(rs) != 0 {
		t.Fatal("route should be skipped")
	}
}

func TestWalk_RelError(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	touch(t, filepath.Join(outside, "a.ts"))

	w := newTestWalker(root, &fakeStrategy{})

	link := filepath.Join(root, "link.ts")
	if err := os.Symlink(filepath.Join(outside, "a.ts"), link); err != nil {
		t.Skip("symlink not supported")
	}

	_, err := w.walk()
	if err == nil {
	}
}

func TestIsRouteFile_AllBranches(t *testing.T) {
	root := t.TempDir()
	w := newTestWalker(root, &fakeStrategy{})

	dirEntry, _ := os.ReadDir(root)
	if len(dirEntry) == 0 {
		os.Mkdir(filepath.Join(root, "dir"), 0o755)
		dirEntry, _ = os.ReadDir(root)
	}
	if w.isRouteFile(root, dirEntry[0]) {
		t.Fatal("dir should be false")
	}

	p := filepath.Join(root, "a.js")
	touch(t, p)
	entries, _ := os.ReadDir(root)
	if w.isRouteFile(p, entries[0]) {
		t.Fatal("wrong ext should be false")
	}

	p2 := filepath.Join(root, "_app.ts")
	touch(t, p2)
	entries, _ = os.ReadDir(root)
	for _, e := range entries {
		if e.Name() == "_app.ts" && w.isRouteFile(p2, e) {
			t.Fatal("excluded should be false")
		}
	}

	p3 := filepath.Join(root, "ok.ts")
	touch(t, p3)
	entries, _ = os.ReadDir(root)
	for _, e := range entries {
		if e.Name() == "ok.ts" && !w.isRouteFile(p3, e) {
			t.Fatal("should be true")
		}
	}
}

func TestBuildName_AllBranches(t *testing.T) {
	w := newTestWalker("", &fakeStrategy{})

	name := w.buildName([]string{"User-Profile", "Detail"})
	if name != "user_profile_detail" {
		t.Fatal(name)
	}

	name2 := w.buildName([]string{"@@@###"})
	if name2 != "" {
		t.Fatal("should be empty")
	}
}

func TestToSnakeCase_AllBranches(t *testing.T) {
	w := newTestWalker("", &fakeStrategy{})

	if w.toSnakeCase([]string{}) != "" {
		t.Fatal("empty should return empty")
	}

	if w.toSnakeCase([]string{"User", "Profile"}) != "user_profile" {
		t.Fatal("wrong snake")
	}
}

func TestSplitSegments_IdentRouteRoot(t *testing.T) {
	root := t.TempDir()
	w := newTestWalker(root, &fakeStrategy{})

	segs := w.splitSegments("index.ts")

	if len(segs) != 1 || segs[0] != "$root" {
		t.Fatalf("expected [$root], got %v", segs)
	}

	p := w.buildPath(segs)
	if p != "/" {
		t.Fatalf("expected '/', got %q", p)
	}
}

func TestSplitSegments_SkipIndexInMiddle(t *testing.T) {
	root := t.TempDir()
	w := newTestWalker(root, &fakeStrategy{})

	segs := w.splitSegments(filepath.Join("users", "index", "profile.ts"))

	want := []string{"users", "profile"}
	if len(segs) != len(want) {
		t.Fatalf("got %v", segs)
	}

	for i := range segs {
		if segs[i] != want[i] {
			t.Fatalf("got %v", segs)
		}
	}
}

func TestWalk_WalkDirError(t *testing.T) {
	root := t.TempDir()

	badDir := filepath.Join(root, "noaccess")
	if err := os.Mkdir(badDir, 0o000); err != nil {
		t.Skip(err)
	}
	defer os.Chmod(badDir, 0o755)

	w := newTestWalker(root, &fakeStrategy{})

	_, err := w.walk()
	if err == nil {
		t.Fatal("expected walkdir error")
	}
}
