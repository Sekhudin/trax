package routes

import (
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/testutil/configmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
)

func TestConfig_Success(t *testing.T) {
	mockConfig := configmock.Config{}

	t.Run("valid_file_strategy", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "file",
				File:     "routes.json",
				Output:   "out.ts",
				Symbols:  &config.RoutesSymbols{},
			}
		}

		cfg, err := NewConfig(mockConfig.Routes()).Load()
		if err != nil || cfg.File == nil || !cfg.IsFileStrategy {
			t.Fatal("file_config_failed")
		}
	})

	t.Run("next_app_normalization", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "next-app",
				Root:     "src",
				Output:   "out.ts",
				Symbols:  &config.RoutesSymbols{},
			}
		}

		cfg, err := NewConfig(mockConfig.Routes()).Load()
		if err != nil || cfg.Root != "src/app" {
			t.Fatal("root_normalization_failed")
		}
	})

	t.Run("next_page_normalization", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "next-page",
				Root:     "src",
				Output:   "out.ts",
				Symbols:  &config.RoutesSymbols{},
			}
		}

		cfg, err := NewConfig(mockConfig.Routes()).Load()
		if err != nil || cfg.Root != "src/pages" {
			t.Fatal("root_normalization_failed")
		}
	})

	t.Run("valid_symbols_normalization", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Symbols: &config.RoutesSymbols{
					Param:    "$p",
					Wildcard: "$w",
					Root:     "p",
				},
			}
		}

		r := &rconfig{
			rule: newConfigRule(),
			cfg:  mockConfig.Routes(),
		}

		sym := r.normalizeSymbols()
		if sym.Param != "$p" || sym.Wildcard != "$w" || sym.Root != "p" {
			t.Fatal("symbol_mapping_failed")
		}
	})
}

func TestConfig_Error(t *testing.T) {
	mockConfig := configmock.Config{}

	t.Run("empty_strategy", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "",
			}
		}

		_, err := NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_err_empty")
		}
	})

	t.Run("invalid_strategy", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "invalid",
			}
		}

		_, err := NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_err_empty")
		}
	})

	t.Run("file_strategy_constraints", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "file",
				File:     "",
			}
		}

		_, err := NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_err_missing_file")
		}

		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "next-app",
				File:     "routes.yaml",
			}
		}

		_, err = NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_err_unset_file")
		}
	})

	t.Run("invalid_extensions", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "next-app",
				Output:   "output.txt",
				Symbols:  &config.RoutesSymbols{},
			}
		}

		_, err := NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_err_ext")
		}

		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "file",
				Output:   "output.ts",
				File:     "invalid.txt",
				Symbols:  &config.RoutesSymbols{},
			}
		}
	})

	t.Run("load_file_parse_error", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "file",
				Root:     "src",
				File:     "json",
				Output:   "output.ts",
				Prefix:   "routes",
				Symbols: &config.RoutesSymbols{
					Param:    "$p",
					Wildcard: "$w",
					Root:     "p",
				},
			}
		}

		_, err := NewConfig(mockConfig.Routes()).Load()
		if err == nil {
			t.Fatal("should_catch_file_parse_error")
		}
	})
}

func TestConfig_Fallback(t *testing.T) {
	mockConfig := configmock.Config{}

	t.Run("symbol_default_values", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Symbols: &config.RoutesSymbols{
					Param:    "foo",
					Wildcard: "foo",
					Root:     "foo",
				},
			}
		}

		r := &rconfig{
			rule: newConfigRule(),
			cfg:  mockConfig.Routes(),
		}

		sym := r.normalizeSymbols()
		if sym.Param != "$param" || sym.Wildcard != "$wildcard" {
			t.Fatal("fallback_symbols_failed")
		}
	})

	t.Run("root_path_idempotency", func(t *testing.T) {
		mock.Reset(&mockConfig)
		mockConfig.RoutesFn = func() *config.RoutesConfig {
			return &config.RoutesConfig{
				Strategy: "next-app",
				Root:     "src/app",
				Output:   "output.ts",
				Prefix:   "routes",
				Symbols: &config.RoutesSymbols{
					Param:    "$p",
					Wildcard: "$w",
					Root:     "p",
				},
			}
		}

		r := &rconfig{
			rule: newConfigRule(),
			cfg:  mockConfig.Routes(),
		}

		if r.normalizeRoot() != "src/app" {
			t.Fatal("should_not_append_double")
		}
	})
}
