package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestAsFlat_Success(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		data     map[string]any
		contains []string
	}{
		{
			name:     "simple_key_value",
			prefix:   "",
			data:     map[string]any{"a": 1, "b": "two"},
			contains: []string{"a = 1", "b = two"},
		},
		{
			name:     "with_custom_prefix",
			prefix:   "config",
			data:     map[string]any{"port": 8080},
			contains: []string{"config.port = 8080"},
		},
		{
			name:     "nested_map_structure",
			prefix:   "",
			data:     map[string]any{"db": map[string]any{"host": "localhost"}},
			contains: []string{"db.host = localhost"},
		},
		{
			name:     "array_index_handling",
			prefix:   "",
			data:     map[string]any{"tags": []any{"go", "cli"}},
			contains: []string{"tags.0 = go", "tags.1 = cli"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			ctx := New(buf, Options{})

			if err := ctx.AsFlat(tt.prefix, tt.data); err != nil {
				t.Fatalf("asflat_failed: %v", err)
			}

			out := buf.String()
			for _, search := range tt.contains {
				if !strings.Contains(out, search) {
					t.Errorf("missing_output: %s", search)
				}
			}
		})
	}
}

func TestAsFlat_Error(t *testing.T) {
	errCtx := New(&errorWriter{}, Options{})

	t.Run("handle_basic_write_failure", func(t *testing.T) {
		err := errCtx.AsFlat("", map[string]any{"key": "val"})
		if err == nil {
			t.Fatal("should_return_write_error")
		}
	})

	t.Run("handle_nested_map_write_failure", func(t *testing.T) {
		data := map[string]any{"parent": map[string]any{"child": 1}}
		if err := errCtx.AsFlat("", data); err == nil {
			t.Fatal("should_return_nested_error")
		}
	})

	t.Run("handle_array_write_failure", func(t *testing.T) {
		data := map[string]any{"list": []any{"item1"}}
		if err := errCtx.AsFlat("", data); err == nil {
			t.Fatal("should_return_array_write_error")
		}
	})
}
