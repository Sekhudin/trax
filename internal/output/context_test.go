package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func newCtx(buf *bytes.Buffer, opt Options) *Context {
	return New(buf, opt)
}

func TestNotify_QuietMode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{Quiet: true})

	ctx.Info("scope", "message")

	if buf.Len() != 0 {
		t.Fatal("expected no output")
	}
}

func TestNotify_JSONMode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{JSON: true, NoColor: true})

	ctx.Success("build", "ok")

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatal(err)
	}

	if out["level"] != "success" ||
		out["scope"] != "build" ||
		out["message"] != "ok" {
		t.Fatal(out)
	}
}

func TestNotify_DebugMode_ForcesJSON(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{Debug: true, NoColor: true})

	ctx.Warn("lint", "warn msg")

	if !strings.Contains(buf.String(), `"level": "warn"`) {
		t.Fatal(buf.String())
	}
}

func TestNotify_TextMode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{NoColor: true})

	ctx.Error("runtime", "boom")

	out := buf.String()

	if !strings.Contains(out, "(runtime)") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "boom") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "✖") {
		t.Fatal(out)
	}
}

func TestIcon_ByLevel(t *testing.T) {
	ctx := newCtx(bytes.NewBuffer(nil), Options{NoColor: true})

	tests := []struct {
		level Level
		icon  string
	}{
		{LevelInfo, "ℹ"},
		{LevelSuccess, "✔"},
		{LevelWarn, "⚠"},
		{LevelError, "✖"},
	}

	for _, tt := range tests {
		out := ctx.icon(tt.level)
		if out != tt.icon {
			t.Fatalf("expected %s got %s", tt.icon, out)
		}
	}
}

func TestColorScope_ByLevel(t *testing.T) {
	ctx := newCtx(bytes.NewBuffer(nil), Options{NoColor: true})

	tests := []struct {
		level Level
	}{
		{LevelInfo},
		{LevelSuccess},
		{LevelWarn},
		{LevelError},
	}

	for _, tt := range tests {
		out := ctx.colorScope(tt.level, "x")
		if out != "x" {
			t.Fatal(out)
		}
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelInfo, "info"},
		{LevelSuccess, "success"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.want {
			t.Fatal(tt.level.String())
		}
	}
}

func TestPublicHelpers(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{NoColor: true})

	ctx.Info("a", "1")
	ctx.Success("b", "2")
	ctx.Warn("c", "3")
	ctx.Error("d", "4")

	out := buf.String()

	if !strings.Contains(out, "ℹ") ||
		!strings.Contains(out, "✔") ||
		!strings.Contains(out, "⚠") ||
		!strings.Contains(out, "✖") {
		t.Fatal(out)
	}
}
