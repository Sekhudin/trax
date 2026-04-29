package routes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sekhudin/trax/internal/path"
)

func rawTest_writeFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	fp := filepath.Join(dir, "routes.yaml")
	if err := os.WriteFile(fp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return fp
}

func rawTest_Cfg(fp string) *Config {
	return &Config{
		IsFileStrategy: true,
		File:           &path.FilePath{Full: fp},
	}
}

func TestRawRouteBuilder_Success(t *testing.T) {
	t.Run("read_valid_file", func(t *testing.T) {
		content := "routes:\n  - name: users\n    path: /users"
		fp := rawTest_writeFile(t, content)
		b := NewRawRouteBuilder(rawTest_Cfg(fp))

		_, err := b.Build()
		if err != nil {
			t.Fatal("should_build_successfully")
		}
	})

	t.Run("read_next_app", func(t *testing.T) {
		cfg := &Config{Strategy: "next-app", Root: "."}
		b := NewRawRouteBuilder(cfg)
		_, _ = b.Build()
	})

	t.Run("read_next_page", func(t *testing.T) {
		cfg := &Config{Strategy: "next-page", Root: "."}
		b := NewRawRouteBuilder(cfg)
		_, _ = b.Build()
	})
}

func TestRawRouteBuilder_Error(t *testing.T) {
	t.Run("config_not_found", func(t *testing.T) {
		b := rawroute{cfg: rawTest_Cfg("/not/exist.yaml")}
		_, err := b.readFile()
		if err == nil {
			t.Fatal("should_trigger_entry_1")
		}
	})

	t.Run("bad_yaml_format", func(t *testing.T) {
		fp := rawTest_writeFile(t, "routes: name 123]")
		b := rawroute{cfg: rawTest_Cfg(fp)}
		_, err := b.readFile()
		if err == nil {
			t.Fatal("should_trigger_entry_2")
		}
	})

	t.Run("empty_routes_config", func(t *testing.T) {
		fp := rawTest_writeFile(t, "routes: []")
		b := rawroute{cfg: rawTest_Cfg(fp)}
		_, err := b.readFile()
		if err == nil {
			t.Fatal("should_trigger_entry_3")
		}
	})

	t.Run("invalid_strategy_type", func(t *testing.T) {
		b := rawroute{cfg: &Config{Strategy: "wrong"}}
		_, err := b.readDisc()
		if err == nil {
			t.Fatal("should_catch_invalid_strategy")
		}
	})

	t.Run("walker_app_failed", func(t *testing.T) {
		cfg := &Config{Strategy: "next-app", Root: "/invalid"}
		b := NewRawRouteBuilder(cfg)
		_, err := b.Build()
		if err == nil {
			t.Fatal("should_catch_walker_error")
		}
	})
}

func TestRawRouteBuilder_Fallback(t *testing.T) {
	t.Run("interface_compliance_check", func(t *testing.T) {
		var _ RawRouteBuilder = (*rawroute)(nil)
	})
}
