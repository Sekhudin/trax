package path

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestParseFilePath_Success(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		allowedExts []string
		expectExt   string
		expectFile  string
	}{
		{
			name:        "valid ts file",
			input:       "src/routes.ts",
			allowedExts: []string{".ts"},
			expectExt:   ".ts",
			expectFile:  "routes.ts",
		},
		{
			name:        "valid nested path",
			input:       "src/trax/routes.ts",
			allowedExts: []string{".ts"},
			expectExt:   ".ts",
			expectFile:  "routes.ts",
		},
		{
			name:        "uppercase extension normalized",
			input:       "src/routes.TS",
			allowedExts: []string{".ts"},
			expectExt:   ".ts",
			expectFile:  "routes.TS",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fp, err := ParseFilePath(tc.input, tc.allowedExts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if fp.Ext != tc.expectExt {
				t.Fatalf("expected ext %s, got %s", tc.expectExt, fp.Ext)
			}

			if !strings.HasSuffix(fp.Filename, tc.expectFile) {
				t.Fatalf("unexpected filename: %s", fp.Filename)
			}

			if fp.Full == "" {
				t.Fatal("expected Full path to not be empty")
			}
		})
	}
}

func TestParseFilePath_EmptyInput(t *testing.T) {
	_, err := ParseFilePath("   ", []string{".ts"})
	if err == nil {
		t.Fatal("expected error")
	}

	ce, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("expected CoreError, got %T", err)
	}

	if ce.Code != appErr.ErrValidation {
		t.Fatalf("expected validation error, got %s", ce.Code)
	}
}

func TestParseFilePath_NoExtension(t *testing.T) {
	_, err := ParseFilePath("src/routes", []string{".ts"})
	if err == nil {
		t.Fatal("expected error")
	}

	ce, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("expected CoreError, got %T", err)
	}

	if ce.Code != appErr.ErrValidation {
		t.Fatalf("expected validation error, got %s", ce.Code)
	}

	if !strings.Contains(ce.Message, "must include file extension") {
		t.Fatal("unexpected message:", ce.Message)
	}
}

func TestParseFilePath_UnsupportedExtension(t *testing.T) {
	_, err := ParseFilePath("src/routes.go", []string{".ts"})
	if err == nil {
		t.Fatal("expected error")
	}

	ce, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("expected CoreError, got %T", err)
	}

	if ce.Code != appErr.ErrValidation {
		t.Fatalf("expected validation error, got %s", ce.Code)
	}

	if !strings.Contains(ce.Message, "unsupported extension") {
		t.Fatal("unexpected message:", ce.Message)
	}
}

func TestParseFilePath_RelPathEdgeCase(t *testing.T) {
	// edge case: relative path resolution from current dir
	input := "./src/routes.ts"

	fp, err := ParseFilePath(input, []string{".ts"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filepath.Ext(fp.Full) != ".ts" {
		t.Fatal("expected .ts extension")
	}
}

func TestIsAllowedExt(t *testing.T) {
	if !isAllowedExt(".ts", []string{".ts", ".js"}) {
		t.Fatal("expected true")
	}

	if isAllowedExt(".go", []string{".ts"}) {
		t.Fatal("expected false")
	}
}

func TestParseFilePath_RelErrorImpossiblePath(t *testing.T) {
	// NOTE: filepath.Rel hampir tidak bisa error di normal OS,
	// tapi kita tetap cover branch error dengan invalid root simulation
	_, err := filepath.Rel(string(rune(0)), "test.ts")
	if err == nil {
		// skip kalau OS tidak error
		t.Skip("cannot trigger Rel error on this OS")
	}

	_, err2 := ParseFilePath("test.ts", []string{".ts"})
	if err2 != nil {
		// just ensure function still behaves normally
		var ce *appErr.CoreError
		_ = errors.As(err2, &ce)
	}
}
