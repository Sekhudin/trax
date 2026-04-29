package config

import (
	"reflect"
	"testing"
)

func TestNew_Success(t *testing.T) {
	t.Run("return_valid_instance", func(t *testing.T) {
		cfg := New()
		if cfg == nil {
			t.Fatal("expected_non_nil")
		}
	})

	t.Run("implement_config_interface", func(t *testing.T) {
		var _ Config = New()
	})

	t.Run("return_concrete_type", func(t *testing.T) {
		cfg := New()
		expected := "*config.config"
		actual := reflect.TypeOf(cfg).String()

		if actual != expected {
			t.Fatalf("expected_%s_got_%s", expected, actual)
		}
	})
}
