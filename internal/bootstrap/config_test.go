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

func TestLoadConfig_ConfigFileProvided(t *testing.T) {
	t.Run("should load successfully when valid config file is provided", func(t *testing.T) {
		resetViper()

		fp := writeTempConfig(t, `
formatter: prettier
routes:
  strategy: custom
`)

		err := LoadConfig(fp)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should return wrapped error when config file is invalid", func(t *testing.T) {
		resetViper()

		fp := writeTempInvalidConfig(t)

		err := LoadConfig(fp)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("should error when config file path does not exist", func(t *testing.T) {
		resetViper()

		err := LoadConfig("/path/does/not/exist.yaml")
		if err == nil {
			t.Fatal("expected error for not found config file, ")
		}
	})
}

func TestLoadConfig_DefaultSearchPath(t *testing.T) {
	t.Run("should not error when no config file found in default path", func(t *testing.T) {
		resetViper()

		dir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		os.Chdir(dir)

		err := LoadConfig("")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("should load config when trax.yaml exists in current directory", func(t *testing.T) {
		resetViper()

		dir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		os.Chdir(dir)

		if err := os.WriteFile(filepath.Join(dir, "trax.yaml"), []byte(`
formatter: prettier
`), 0o644); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		err := LoadConfig("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should return error when trax.yaml is invalid", func(t *testing.T) {
		resetViper()

		dir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		os.Chdir(dir)

		if err := os.WriteFile(filepath.Join(dir, "trax.yaml"), []byte("invalid: [yaml"), 0o644); err != nil {
			t.Fatalf("failed to write invalid config: %v", err)
		}

		err := LoadConfig("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
