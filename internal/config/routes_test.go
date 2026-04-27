package config

import (
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func TestRoutes_DefaultValues(t *testing.T) {
	resetViper()

	cfg := New()
	routes := cfg.Routes()

	if routes == nil {
		t.Fatal("expected routes config not nil")
	}

	if routes.Symbols == nil {
		t.Fatal("expected symbols not nil")
	}

	if routes.Strategy != "" ||
		routes.Root != "" ||
		routes.File != "" ||
		routes.Output != "" ||
		routes.Prefix != "" ||
		routes.NoDeps != false ||
		routes.Symbols.Param != "" ||
		routes.Symbols.Wildcard != "" ||
		routes.Symbols.Root != "" {
		t.Fatal("expected all default zero values")
	}
}

func TestRoutes_ReadsAllValuesFromViper(t *testing.T) {
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

	// top level
	if routes.Strategy != "tree" {
		t.Fatal("strategy not read correctly")
	}
	if routes.Root != "/app" {
		t.Fatal("root not read correctly")
	}
	if routes.NoDeps != true {
		t.Fatal("no-deps not read correctly")
	}
	if routes.File != "routes.go" {
		t.Fatal("file not read correctly")
	}
	if routes.Output != "generated" {
		t.Fatal("output not read correctly")
	}
	if routes.Prefix != "/api" {
		t.Fatal("prefix not read correctly")
	}

	// symbols
	if routes.Symbols.Param != ":" {
		t.Fatal("param symbol not read correctly")
	}
	if routes.Symbols.Wildcard != "*" {
		t.Fatal("wildcard symbol not read correctly")
	}
	if routes.Symbols.Root != "/" {
		t.Fatal("root symbol not read correctly")
	}
}
