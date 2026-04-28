package routes

import (
	"testing"

	"github.com/sekhudin/trax/internal/config"
)

func TestNewConfig_StrategyEmpty(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{}).Load()
	if err == nil {
		t.Fatal("expected error for empty strategy")
	}
}

func TestNewConfig_InvalidStrategy(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{
		Strategy: "invalid",
	}).Load()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestNewConfig_FileStrategy_MissingFile(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{
		Strategy: "file",
	}).Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewConfig_NonFileStrategy_WithFileSet(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{
		Strategy: "next-app",
		File:     "routes.json",
	}).Load()
	if err == nil {
		t.Fatal("expected error when file is set for non-file strategy")
	}
}

func TestNewConfig_FileStrategy_InvalidFileExt(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{
		Strategy: "file",
		Root:     "src",
		File:     "routes.txt",
		Output:   "out.ts",
		Prefix:   "routes",
		NoDeps:   false,
		Symbols: &config.RoutesSymbols{
			Param:    "$p",
			Wildcard: "$w",
			Root:     "p",
		},
	}).Load()
	if err == nil {
		t.Fatal("expected extension error")
	}
}

func TestNewConfig_OutputInvalidExt(t *testing.T) {
	_, err := NewConfig(&config.RoutesConfig{
		Strategy: "file",
		File:     "routes.json",
		Output:   "out.txt",
	}).Load()
	if err == nil {
		t.Fatal("expected output extension error")
	}
}

func TestNewConfig_ValidFileStrategy(t *testing.T) {
	cfg, err := NewConfig(&config.RoutesConfig{
		Strategy: "file",
		File:     "routes.json",
		Output:   "out.ts",
		Symbols:  &config.RoutesSymbols{},
	}).Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.File == nil || cfg.Output == nil {
		t.Fatal("file or output not parsed")
	}

	if !cfg.IsFileStrategy {
		t.Fatal("expected file strategy true")
	}
}

func TestNewConfig_ValidNextAppStrategy(t *testing.T) {
	cfg, err := NewConfig(&config.RoutesConfig{
		Strategy: "next-app",
		Root:     "src",
		Output:   "out.ts",
		Symbols:  &config.RoutesSymbols{},
	}).Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Root != "src/app" {
		t.Fatalf("expected src/app, got %s", cfg.Root)
	}
}

func TestNewConfig_ValidNextPageStrategy(t *testing.T) {
	cfg, err := NewConfig(&config.RoutesConfig{
		Strategy: "next-page",
		Output:   "out.ts",
		Root:     "src",
		Symbols:  &config.RoutesSymbols{},
	}).Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Root != "src/pages" {
		t.Fatalf("expected src/pages, got %s", cfg.Root)
	}
}

func TestNormalizeSymbols_DefaultFallback(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Symbols: &config.RoutesSymbols{
				Param:    "invalid",
				Wildcard: "invalid",
				Root:     "invalid",
			},
		},
	}

	sym := r.normalizeSymbols()

	if sym.Param != "$param" ||
		sym.Wildcard != "$wildcard" ||
		sym.Root != "root" {
		t.Fatal(sym)
	}
}

func TestNormalizeSymbols_ValidSymbols(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Symbols: &config.RoutesSymbols{
				Param:    "$p",
				Wildcard: "$w",
				Root:     "p",
			},
		},
	}

	sym := r.normalizeSymbols()

	if sym.Param != "$p" ||
		sym.Wildcard != "$w" ||
		sym.Root != "p" {
		t.Fatal(sym)
	}
}

func TestNormalizeRootPath_NoSuffix(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Strategy: "file",
			Root:     "src",
		},
	}

	out := r.normalizeRoot()

	if out != "src" {
		t.Fatal(out)
	}
}

func TestNormalizeRootPath_AlreadyHasSuffix(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Strategy: "next-app",
			Root:     "src/app",
		},
	}

	out := r.normalizeRoot()

	if out != "src/app" {
		t.Fatal(out)
	}
}

func TestNormalizeRootPath_AddSuffix(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Strategy: "next-page",
			Root:     "src",
		},
	}

	out := r.normalizeRoot()

	if out != "src/pages" {
		t.Fatal(out)
	}
}

func TestIsFileStrategy(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Strategy: "file",
		},
	}

	if !r.IsFileStrategy() {
		t.Fatal("should be file strategy")
	}

	r.cfg.Strategy = "next-app"
	if r.IsFileStrategy() {
		t.Fatal("should not be file strategy")
	}
}

func TestIsValidStrategy(t *testing.T) {
	r := &rconfig{
		rule: newConfigRule(),
		cfg: &config.RoutesConfig{
			Strategy: "file",
		},
	}

	if !r.IsValidStrategy() {
		t.Fatal("valid strategy not detected")
	}

	r.cfg.Strategy = "unknown"
	if r.IsValidStrategy() {
		t.Fatal("invalid strategy detected as valid")
	}
}
