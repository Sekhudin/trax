package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestAsJSON_Success(t *testing.T) {
	t.Run("complex_data_structure", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := New(buf, Options{})

		data := map[string]any{
			"a": 1,
			"b": map[string]any{"c": 2},
			"d": []any{3, "e"},
		}

		if err := ctx.AsJSON(data); err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(buf.String(), "\n") || !strings.Contains(buf.String(), `"c": 2`) {
			t.Fatal("format_mismatch")
		}
	})
}

func TestAsJSON_Error(t *testing.T) {
	t.Run("stdout_write_failed", func(t *testing.T) {
		ctx := New(&errorWriter{}, Options{})
		err := ctx.AsJSON(map[string]any{"a": 1})

		if err == nil {
			t.Fatal("should_error")
		}

		coreErr, ok := err.(*appErr.CoreError)
		if !ok || coreErr.Code != appErr.ErrIO {
			t.Fatalf("wrong_error_code: %v", err)
		}
	})

	t.Run("json_marshal_failed", func(t *testing.T) {
		ctx := New(bytes.NewBuffer(nil), Options{
			Marshal: func(v any, prefix, indent string) ([]byte, error) {
				return nil, errors.New("boom")
			},
		})

		err := ctx.AsJSON(map[string]any{"a": 1})
		if err == nil {
			t.Fatal("should_error")
		}

		coreErr, ok := err.(*appErr.CoreError)
		if !ok || coreErr.Code != appErr.ErrInternal {
			t.Fatalf("wrong_error_code: %v", err)
		}
	})
}

func TestAsJSON_Fallback(t *testing.T) {
	t.Run("normalization_recursion_logic", func(t *testing.T) {
		ctx := New(bytes.NewBuffer(nil), Options{}).(*context)

		in := map[string]any{"p": map[string]any{"c": 1}}
		out := ctx.normalizeValue(in).(map[string]any)

		if _, ok := out["p"].(map[string]any); !ok {
			t.Fatal("recursion_failed")
		}
	})
}
