package routes

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type mockStrategy struct {
	skipErr error
}

func (f *mockStrategy) shouldSkip(string, fs.DirEntry) error { return f.skipErr }
func (f *mockStrategy) normalizeSegment(seg string) string {
	return seg
}

func newTestWalker(root string, strat walkerstrategy) *walker {
	return &walker{
		strategy: strat,
		config:   &Config{Root: root},
		rule: &walkrule{
			exts:         map[string]struct{}{".ts": {}},
			identRoute:   map[string]struct{}{"index": {}},
			excludeFiles: map[string]struct{}{"_app": {}},
		},
	}
}

func testWalker_Touch(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestWalker_Success(t *testing.T) {
	root := t.TempDir()
	w := newTestWalker(root, &mockStrategy{})

	t.Run("walk_valid_routes", func(t *testing.T) {
		testWalker_Touch(t, filepath.Join(root, "users", "index.ts"))
		testWalker_Touch(t, filepath.Join(root, "users", "profile.ts"))
		testWalker_Touch(t, filepath.Join(root, "_app.ts"))
		testWalker_Touch(t, filepath.Join(root, "ignore.js"))

		rs, err := w.walk()
		if err != nil || len(rs) != 2 {
			t.Fatal("fail")
		}
	})

	t.Run("is_route_file", func(t *testing.T) {
		p := filepath.Join(root, "ok.ts")
		testWalker_Touch(t, p)
		entries, _ := os.ReadDir(root)
		for _, e := range entries {
			if e.Name() == "ok.ts" && !w.isRouteFile(p, e) {
				t.Fatal("fail")
			}
		}
	})

	t.Run("build_snake_name", func(t *testing.T) {
		name := w.buildName([]string{"User-Profile", "Detail"})
		if name != "user_profile_detail" {
			t.Fatal("fail")
		}
	})

	t.Run("to_snake_case", func(t *testing.T) {
		if w.toSnakeCase([]string{"User", "Profile"}) != "user_profile" {
			t.Fatal("fail")
		}
	})

	t.Run("split_root_index", func(t *testing.T) {
		segs := w.splitSegments("index.ts")
		if len(segs) != 1 || segs[0] != "$root" {
			t.Fatal("fail")
		}
		if w.buildPath(segs) != "/" {
			t.Fatal("fail")
		}
	})

	t.Run("skip_middle_index", func(t *testing.T) {
		segs := w.splitSegments(filepath.Join("users", "index", "profile.ts"))
		if len(segs) != 2 || segs[1] != "profile" {
			t.Fatal("fail")
		}
	})
}

func TestWalker_Error(t *testing.T) {
	root := t.TempDir()

	t.Run("strategy_skip_error", func(t *testing.T) {
		testWalker_Touch(t, filepath.Join(root, "a.ts"))
		w := newTestWalker(root, &mockStrategy{skipErr: errors.New("stop")})
		if _, err := w.walk(); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("walk_dir_access", func(t *testing.T) {
		badDir := filepath.Join(root, "noaccess")
		os.Mkdir(badDir, 0o000)
		defer os.Chmod(badDir, 0o755)
		w := newTestWalker(root, &mockStrategy{})
		if _, err := w.walk(); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("relative_path_error", func(t *testing.T) {
		old := filepathRel
		filepathRel = func(base, targ string) (string, error) {
			return "", errors.New("mock_rel_error")
		}

		t.Cleanup(func() { filepathRel = old })

		outside := t.TempDir()
		testWalker_Touch(t, filepath.Join(outside, "a.ts"))
		w := newTestWalker(root, &mockStrategy{})
		link := filepath.Join(root, "link.ts")
		if err := os.Symlink(filepath.Join(outside, "a.ts"), link); err != nil {
			t.Skip("no symlink")
		}
		_, _ = w.walk()
	})
}

func TestWalker_Fallback(t *testing.T) {
	root := t.TempDir()
	w := newTestWalker(root, &mockStrategy{})

	t.Run("skip_directory_entry", func(t *testing.T) {
		os.Mkdir(filepath.Join(root, "dir"), 0o755)
		entries, _ := os.ReadDir(root)
		if w.isRouteFile(root, entries[0]) {
			t.Fatal("fail")
		}
	})

	t.Run("ignore_wrong_extension", func(t *testing.T) {
		p := filepath.Join(root, "a.js")
		testWalker_Touch(t, p)
		entries, _ := os.ReadDir(root)
		if w.isRouteFile(p, entries[0]) {
			t.Fatal("fail")
		}
	})

	t.Run("exclude_app_file", func(t *testing.T) {
		p := filepath.Join(root, "_app.ts")
		testWalker_Touch(t, p)
		entries, _ := os.ReadDir(root)
		for _, e := range entries {
			if e.Name() == "_app.ts" && w.isRouteFile(p, e) {
				t.Fatal("fail")
			}
		}
	})

	t.Run("sanitize_empty_input", func(t *testing.T) {
		if w.buildName([]string{"@@@###"}) != "" {
			t.Fatal("fail")
		}
	})

	t.Run("empty_snake_input", func(t *testing.T) {
		if w.toSnakeCase([]string{}) != "" {
			t.Fatal("fail")
		}
	})
}
