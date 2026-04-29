package config

import (
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func TestRoutes_Success(t *testing.T) {
	t.Run("read_all_values", func(t *testing.T) {
		resetViper()

		viper.Set("routes.strategy", "tree")
		viper.Set("routes.root", "/app")
		viper.Set("routes.no-deps", true)
		viper.Set("routes.file", "routes.go")
		viper.Set("routes.output", "generated")
		viper.Set("routes.prefix", "/api")
		viper.Set("routes.symbols.param", ":")
		viper.Set("routes.symbols.wildcard", "*")
		viper.Set("routes.symbols.root", "/")

		cfg := New()
		routes := cfg.Routes()

		if routes.Strategy != "tree" || routes.Root != "/app" || !routes.NoDeps {
			t.Fatal("failed_read_toplevel")
		}
		if routes.File != "routes.go" || routes.Output != "generated" || routes.Prefix != "/api" {
			t.Fatal("failed_read_paths")
		}

		if routes.Symbols.Param != ":" || routes.Symbols.Wildcard != "*" || routes.Symbols.Root != "/" {
			t.Fatal("failed_read_symbols")
		}
	})
}

func TestRoutes_Fallback(t *testing.T) {
	t.Run("return_zero_values", func(t *testing.T) {
		resetViper()

		cfg := New()
		routes := cfg.Routes()

		if routes == nil {
			t.Fatal("expected_non_nil")
		}

		if routes.Symbols == nil {
			t.Fatal("symbols_is_nil")
		}

		if routes.Strategy != "" || routes.NoDeps != false || routes.Symbols.Param != "" {
			t.Fatal("expected_defaults")
		}
	})
}
