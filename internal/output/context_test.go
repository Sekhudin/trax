package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced write error")
}

func TestOutput_Success(t *testing.T) {
	t.Run("format_text_output", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := New(buf, Options{NoColor: true})

		ctx.Error("runtime", "boom")
		ctx.Cause("root", "denied")

		out := buf.String()
		if !strings.Contains(out, "✖") || !strings.Contains(out, "↳") {
			t.Fatal("wrong_text_format")
		}
	})

	t.Run("format_json_output", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := New(buf, Options{JSON: true})

		ctx.Success("build", "ok")

		var out map[string]any
		if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
			t.Fatal("invalid_json_output")
		}

		if out["level"] != "success" || out["scope"] != "build" {
			t.Fatal("wrong_json_fields")
		}
	})

	t.Run("debug_mode_forces_json", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := New(buf, Options{Debug: true})

		ctx.Warn("lint", "msg")
		if !strings.Contains(buf.String(), `"level": "warn"`) {
			t.Fatal("debug_not_json")
		}
	})
}

func TestOutput_Fallback(t *testing.T) {
	t.Run("handle_quiet_mode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := New(buf, Options{Quiet: true})

		ctx.Info("any", "msg")
		if buf.Len() != 0 {
			t.Fatal("output_not_muted")
		}
	})

	t.Run("handle_marshal_error", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		ctx := New(buf, Options{
			JSON: true,
			Marshal: func(v any, prefix, indent string) ([]byte, error) {
				return nil, fmt.Errorf("forced error")
			},
		})

		ctx.Info("test", "msg")

		if !strings.Contains(buf.String(), "failed to marshal json") {
			t.Fatal("fallback_message_not_found")
		}
	})
}

func TestLevel_Internal(t *testing.T) {
	tests := []struct {
		level Level
		icon  string
		str   string
	}{
		{LevelInfo, "ℹ", "info"},
		{LevelSuccess, "✔", "success"},
		{LevelWarn, "⚠", "warn"},
		{LevelError, "✖", "error"},
		{LevelCause, "↳", "cause"},
	}

	ctx := &context{color: NewColorizer(true)}

	for _, tt := range tests {
		t.Run("variant_"+tt.str, func(t *testing.T) {
			if strings.TrimSpace(ctx.icon(tt.level)) != strings.TrimSpace(tt.icon) {
				t.Errorf("icon_mismatch_%s", tt.str)
			}

			if ctx.colorScope(tt.level, "test") != "test" {
				t.Errorf("scope_mismatch_%s", tt.str)
			}

			if tt.level.String() != tt.str {
				t.Errorf("string_mismatch_%s", tt.str)
			}
		})
	}

	t.Run("check_unknown_fallback", func(t *testing.T) {
		unknown := Level(99)
		if unknown.String() != "info" {
			t.Fatal("wrong_default_string")
		}
		if strings.TrimSpace(ctx.icon(unknown)) != "ℹ" {
			t.Fatal("wrong_default_icon")
		}
	})
}
