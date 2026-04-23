package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func reset(t *testing.T) {
	t.Helper()
	viper.Reset()
}

func chdir(t *testing.T, dir string) {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(old)
	})
}

func TestLoadConfig_Success(t *testing.T) {
	reset(t)
	t.Cleanup(func() { viper.Reset() })

	dir := t.TempDir()
	chdir(t, dir)

	cfgPath := filepath.Join(dir, "trax.yaml")

	err := os.WriteFile(cfgPath, []byte(strings.Join([]string{
		"debug: true",
		"user:",
		` name: "john doe"`,
	}, "\n")), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	err = LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !viper.GetBool("debug") {
		t.Fatalf("expected debug true")
	}

	if viper.GetString("user.name") != "john doe" {
		t.Fatalf("expected user.name from file")
	}
}

func TestLoadConfig_NotFound(t *testing.T) {
	reset(t)
	t.Cleanup(func() { viper.Reset() })

	err := LoadConfig("file-yang-tidak-ada.yaml")
	if err == nil {
		t.Fatal("expected error")
	}

	coreErr, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("expected CoreError")
	}

	if coreErr.Code != appErr.ErrConfigNotFound {
		t.Log("CODE", coreErr.Code)
		t.Fatalf("expected ErrConfigNotFound")
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	reset(t)
	t.Cleanup(func() { viper.Reset() })

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "trax.ts")

	err := os.WriteFile(cfgPath, []byte(`:::::::`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	err = LoadConfig(cfgPath)
	if err == nil {
		t.Fatal("expected error")
	}

	coreErr, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("expected CoreError")
	}

	if coreErr.Code != appErr.ErrConfigLoad {
		t.Fatalf("expected ErrConfigLoad")
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	reset(t)
	t.Cleanup(func() { viper.Reset() })

	os.Setenv("TRAX_USER_NAME", "override")
	t.Cleanup(func() { os.Unsetenv("TRAX_USER_NAME") })

	err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	name := viper.GetString("user.name")
	if name != "override" {
		t.Fatalf("expected env override, got %s", name)
	}
}
