package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	fp := filepath.Join(dir, "trax.yaml")

	if err := os.WriteFile(fp, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	return fp
}

func writeTempInvalidConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	fp := filepath.Join(dir, "trax.yaml")

	if err := os.WriteFile(fp, []byte("formatter: [invalid_yaml"), 0o644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	return fp
}

func TestLoadConfig_Success(t *testing.T) {
	t.Run("load_provided_file", func(t *testing.T) {
		resetViper()
		fp := writeTempConfig(t, "formatter: prettier")

		if err := LoadConfig(fp); err != nil {
			t.Fatalf("expected_no_error: %v", err)
		}
	})

	t.Run("load_default_file", func(t *testing.T) {
		resetViper()
		dir := t.TempDir()

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(dir)

		content := []byte("formatter: prettier")
		os.WriteFile(filepath.Join(dir, "trax.yaml"), content, 0o644)

		if err := LoadConfig(""); err != nil {
			t.Fatalf("expected_no_error: %v", err)
		}
	})

	t.Run("handle_missing_default", func(t *testing.T) {
		resetViper()
		if err := LoadConfig(""); err != nil {
			t.Fatalf("expected_nil_error: %v", err)
		}
	})
}

func TestLoadConfig_Error(t *testing.T) {
	t.Run("invalid_yaml_format", func(t *testing.T) {
		resetViper()
		fp := writeTempInvalidConfig(t)

		if err := LoadConfig(fp); err == nil {
			t.Fatal("expected_error_invalid_yaml")
		}
	})

	t.Run("missing_custom_path", func(t *testing.T) {
		resetViper()
		if err := LoadConfig("/non/existent/trax.yaml"); err == nil {
			t.Fatal("expected_error_not_found")
		}
	})

	t.Run("invalid_default_content", func(t *testing.T) {
		resetViper()
		dir := t.TempDir()

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(dir)

		os.WriteFile(filepath.Join(dir, "trax.yaml"), []byte("invalid: ["), 0o644)

		if err := LoadConfig(""); err == nil {
			t.Fatal("expected_error_corrupt_file")
		}
	})
}
