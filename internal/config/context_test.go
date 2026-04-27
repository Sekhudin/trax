package config

import (
	"reflect"
	"testing"
)

func TestNew_ReturnsConfigImplementation(t *testing.T) {
	cfg := New()

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestNew_ImplementsConfigInterface(t *testing.T) {
	var cfg Config = New()

	if cfg == nil {
		t.Fatal("expected config to implement Config interface")
	}
}

func TestNew_ReturnsConcreteTypeConfig(t *testing.T) {
	cfg := New()

	expectedType := "*config.config"
	actualType := reflect.TypeOf(cfg).String()

	if actualType != expectedType {
		t.Fatalf("expected type %s, got %s", expectedType, actualType)
	}
}
