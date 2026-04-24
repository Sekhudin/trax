package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
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
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func TestLoadConfig_Success(t *testing.T) {
	reset(t)

	dir := t.TempDir()
	chdir(t, dir)

	cfgPath := filepath.Join(dir, "trax.yaml")
	err := os.WriteFile(cfgPath, []byte(strings.Join([]string{
		"formatter: prettier",
		"routes:",
		"  prefix: custom-routes",
	}, "\n")), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	err = LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if viper.GetString("formatter") != "prettier" {
		t.Fatalf("expected formatter prettier from file, got %s", viper.GetString("formatter"))
	}

	if viper.GetString("routes.prefix") != "custom-routes" {
		t.Fatalf("expected custom prefix, got %s", viper.GetString("routes.prefix"))
	}
}

func TestLoadConfig_FileNotFound_ShouldUseDefault(t *testing.T) {
	reset(t)

	err := LoadConfig("notfound.yaml")
	if err == nil {
		t.Fatalf("expected error when file is missing. got: %v", err)
	}

	if viper.GetString("formatter") != "biome" {
		t.Fatalf("expected default formatter 'biome', got %s", viper.GetString("formatter"))
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	reset(t)

	os.Setenv("TRAX_FORMATTER", "custom-env")
	t.Cleanup(func() { os.Unsetenv("TRAX_FORMATTER") })

	err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if viper.GetString("formatter") != "custom-env" {
		t.Fatalf("expected env override 'custom-env', got %s", viper.GetString("formatter"))
	}
}
