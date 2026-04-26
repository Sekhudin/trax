package routes

import (
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func TestNewConfig_StrategyEmpty(t *testing.T) {
	resetViper()

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected error for empty strategy")
	}
}

func TestNewConfig_InvalidStrategy(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "invalid")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestNewConfig_FileStrategy_MissingFile(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "file")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewConfig_NonFileStrategy_WithFileSet(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "next-app")
	viper.Set("routes.file", "routes.json")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected error when file is set for non-file strategy")
	}
}

func TestNewConfig_FileStrategy_InvalidFileExt(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "file")
	viper.Set("routes.file", "routes.txt")
	viper.Set("routes.output", "out.ts")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected extension error")
	}
}

func TestNewConfig_OutputInvalidExt(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "file")
	viper.Set("routes.file", "routes.json")
	viper.Set("routes.output", "out.txt")

	_, err := NewConfig()
	if err == nil {
		t.Fatal("expected output extension error")
	}
}

func TestNewConfig_ValidFileStrategy(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "file")
	viper.Set("routes.file", "routes.json")
	viper.Set("routes.output", "out.ts")

	cfg, err := NewConfig()
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
	resetViper()

	viper.Set("routes.strategy", "next-app")
	viper.Set("routes.output", "out.ts")
	viper.Set("routes.root", "src")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Root != "src/app" {
		t.Fatalf("expected src/app, got %s", cfg.Root)
	}
}

func TestNewConfig_ValidNextPageStrategy(t *testing.T) {
	resetViper()

	viper.Set("routes.strategy", "next-page")
	viper.Set("routes.output", "out.ts")
	viper.Set("routes.root", "src")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Root != "src/pages" {
		t.Fatalf("expected src/pages, got %s", cfg.Root)
	}
}

func TestNormalizeSymbols_DefaultFallback(t *testing.T) {
	resetViper()

	r := newConfigRule()

	viper.Set("routes.symbols.param", "invalid")
	viper.Set("routes.symbols.wildcard", "invalid")
	viper.Set("routes.symbols.root", "invalid")

	sym := r.normalizeSymbols()

	if sym.Param != "$param" ||
		sym.Wildcard != "$wildcard" ||
		sym.Root != "root" {
		t.Fatal(sym)
	}
}

func TestNormalizeSymbols_ValidSymbols(t *testing.T) {
	resetViper()

	r := newConfigRule()

	viper.Set("routes.symbols.param", "$p")
	viper.Set("routes.symbols.wildcard", "$w")
	viper.Set("routes.symbols.root", "p")

	sym := r.normalizeSymbols()

	if sym.Param != "$p" ||
		sym.Wildcard != "$w" ||
		sym.Root != "p" {
		t.Fatal(sym)
	}
}

func TestNormalizeRootPath_NoSuffix(t *testing.T) {
	r := newConfigRule()

	out := r.normalizeRootPath("src", "file")

	if out != "src" {
		t.Fatal(out)
	}
}

func TestNormalizeRootPath_AlreadyHasSuffix(t *testing.T) {
	r := newConfigRule()

	out := r.normalizeRootPath("src/app", "next-app")

	if out != "src/app" {
		t.Fatal(out)
	}
}

func TestNormalizeRootPath_AddSuffix(t *testing.T) {
	r := newConfigRule()

	out := r.normalizeRootPath("src", "next-page")

	if out != "src/pages" {
		t.Fatal(out)
	}
}

func TestIsFileStrategy(t *testing.T) {
	r := newConfigRule()

	if !r.IsFileStrategy("file") {
		t.Fatal("should be file strategy")
	}

	if r.IsFileStrategy("next-app") {
		t.Fatal("should not be file strategy")
	}
}

func TestIsValidStrategy(t *testing.T) {
	r := newConfigRule()

	if !r.isValidStartegy("file") {
		t.Fatal("valid strategy not detected")
	}

	if r.isValidStartegy("unknown") {
		t.Fatal("invalid strategy detected as valid")
	}
}
