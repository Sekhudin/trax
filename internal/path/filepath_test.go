package path

import (
	"path/filepath"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestParseFilePath_ValidationErrors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := ParseFilePath("   ", []string{".ts"})
		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := err.(*appErr.CoreError); !ok {
			t.Fatal("expected CoreError")
		}
	})

	t.Run("missing extension", func(t *testing.T) {
		_, err := ParseFilePath("file", []string{".ts"})
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "must include file extension") {
			t.Fatal(err)
		}
	})

	t.Run("unsupported extension", func(t *testing.T) {
		_, err := ParseFilePath("file.go", []string{".ts"})
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "unsupported extension") {
			t.Fatal(err)
		}
	})
}

func TestParseFilePath_ValidCases(t *testing.T) {
	t.Run("simple valid file", func(t *testing.T) {
		fp, err := ParseFilePath("test.ts", []string{".ts", ".js"})
		if err != nil {
			t.Fatal(err)
		}

		if fp.Ext != ".ts" {
			t.Fatalf("unexpected ext: %s", fp.Ext)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		fp, err := ParseFilePath("./src/file.js", []string{".js"})
		if err != nil {
			t.Fatal(err)
		}

		if fp.Ext != ".js" {
			t.Fatal("wrong ext")
		}

		if fp.Dir == "" {
			t.Fatal("expected dir")
		}
	})

	t.Run("absolute path normalization", func(t *testing.T) {
		tmp := filepath.Join(".", "example.ts")

		fp, err := ParseFilePath(tmp, []string{".ts"})
		if err != nil {
			t.Fatal(err)
		}

		if !filepath.IsAbs(fp.Full) {
			t.Fatal("expected absolute path")
		}
	})

	t.Run("case insensitive extension", func(t *testing.T) {
		fp, err := ParseFilePath("FILE.TS", []string{".ts"})
		if err != nil {
			t.Fatal(err)
		}

		if fp.Ext != ".ts" {
			t.Fatal("expected lowercase normalized ext")
		}
	})
}

func TestParseFilePath_ExtensionRules(t *testing.T) {
	t.Run("allowed ext single", func(t *testing.T) {
		_, err := ParseFilePath("a.ts", []string{".ts"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("allowed ext multiple", func(t *testing.T) {
		_, err := ParseFilePath("a.js", []string{".ts", ".js"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("not allowed ext", func(t *testing.T) {
		_, err := ParseFilePath("a.go", []string{".ts", ".js"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty allowed list", func(t *testing.T) {
		_, err := ParseFilePath("a.ts", []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestIsAllowedExt(t *testing.T) {
	t.Run("true case", func(t *testing.T) {
		if !isAllowedExt(".ts", []string{".ts"}) {
			t.Fatal("expected true")
		}
	})

	t.Run("false case", func(t *testing.T) {
		if isAllowedExt(".go", []string{".ts"}) {
			t.Fatal("expected false")
		}
	})

	t.Run("nil allowed list", func(t *testing.T) {
		if isAllowedExt(".ts", nil) {
			t.Fatal("expected false")
		}
	})
}
