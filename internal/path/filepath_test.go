package path

import (
	"errors"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestParseFilePath_Success(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		exts   []string
		expect string
	}{
		{"standard_file_path", "src/routes.ts", []string{".ts"}, ".ts"},
		{"nested_folder_path", "app/api/user.go", []string{".go"}, ".go"},
		{"case_insensitive_ext", "README.MD", []string{".md"}, ".md"},
		{"dot_slash_prefix", "./main.go", []string{".go"}, ".go"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fp, err := ParseFilePath(tc.input, tc.exts)
			if err != nil {
				t.Fatalf("unexpected_error: %v", err)
			}
			if fp.Ext != tc.expect {
				t.Errorf("ext_mismatch: %s", fp.Ext)
			}
		})
	}
}

func TestParseFilePath_Error(t *testing.T) {
	t.Run("empty_path_input", func(t *testing.T) {
		_, err := ParseFilePath("  ", []string{".ts"})
		if err == nil || err.(*appErr.CoreError).Code != appErr.ErrValidation {
			t.Fatal("should_validate_empty")
		}
	})

	t.Run("missing_file_extension", func(t *testing.T) {
		_, err := ParseFilePath("Dockerfile", []string{".ts"})
		if err == nil || !strings.Contains(err.Error(), "must include") {
			t.Fatal("should_require_extension")
		}
	})

	t.Run("unsupported_file_extension", func(t *testing.T) {
		_, err := ParseFilePath("index.php", []string{".ts", ".go"})
		if err == nil || !strings.Contains(err.Error(), "unsupported") {
			t.Fatal("should_block_extension")
		}
	})

	t.Run("rel_path_failed", func(t *testing.T) {
		old := filepathRel
		filepathRel = func(base, targ string) (string, error) {
			return "", errors.New("mock_rel_error")
		}

		t.Cleanup(func() { filepathRel = old })

		_, err := ParseFilePath("test.ts", []string{".ts"})
		if err == nil || err.Error() != "mock_rel_error" {
			t.Fatal("should_catch_rel_error")
		}
	})
}

func TestParseFilePath_Fallback(t *testing.T) {
	t.Run("allowed_ext_checker", func(t *testing.T) {
		if !isAllowedExt(".ts", []string{".ts"}) {
			t.Fatal("check_failed")
		}
	})

	t.Run("path_component_split", func(t *testing.T) {
		fp, _ := ParseFilePath("dir/file.ts", []string{".ts"})
		if fp.Dir != "dir" || fp.Filename != "file.ts" {
			t.Fatal("split_logic_failed")
		}
	})
}
