package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func newCtx(buf *bytes.Buffer, opt Options) Context {
	return New(buf, opt)
}

func newCtxStruct(w *bytes.Buffer, opt Options) *context {
	return &context{
		w:     w,
		opt:   opt,
		color: NewColorizer(opt.NoColor),
	}
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
	ctx := newCtxStruct(bytes.NewBuffer(nil), Options{NoColor: true})

	tests := []struct {
		level Level
		icon  string
	}{
		{LevelInfo, "ℹ"},
		{LevelSuccess, "✔"},
		{LevelWarn, "⚠"},
		{LevelError, "✖"},
		{LevelCause, "   ↳"},
	}

	for _, tt := range tests {
		out := ctx.icon(tt.level)
		if out != tt.icon {
			t.Fatalf("expected %s got %s", tt.icon, out)
		}
	}
}

func TestColorScope_ByLevel(t *testing.T) {
	ctx := newCtxStruct(bytes.NewBuffer(nil), Options{NoColor: true})

	tests := []struct {
		level Level
	}{
		{LevelInfo},
		{LevelSuccess},
		{LevelWarn},
		{LevelError},
		{LevelCause},
	}

	for _, tt := range tests {
		out := ctx.colorScope(tt.level, "x")
		if out != "x" {
			t.Fatal(out)
		}
	}
}

func TestColorContext(t *testing.T) {
	ctx := newCtx(bytes.NewBuffer(nil), Options{NoColor: true})

	if ctx.Color() == nil {
		t.Fatal("expected color not nil")
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
		{LevelCause, "cause"},
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
	ctx.Cause("e", "5")

	out := buf.String()

	if !strings.Contains(out, "ℹ") ||
		!strings.Contains(out, "✔") ||
		!strings.Contains(out, "⚠") ||
		!strings.Contains(out, "✖") ||
		!strings.Contains(out, "   ↳") {
		t.Fatal(out)
	}
}

func TestNotifyJSON_SuccessPath_Explicit(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtx(buf, Options{
		JSON:    true,
		NoColor: true,
	})

	ctx.Info("scope-x", "message-x")

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatal(err)
	}

	if out["level"] != "info" {
		t.Fatalf("level mismatch: %v", out)
	}
	if out["scope"] != "scope-x" {
		t.Fatalf("scope mismatch: %v", out)
	}
	if out["message"] != "message-x" {
		t.Fatalf("message mismatch: %v", out)
	}
}

func TestNotifyJSON_MarshalErrorPath_Fixed(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newCtxStruct(buf, Options{
		JSON:    true,
		NoColor: true,
	})

	bad := string([]byte{0xff, 0xfe, 0xfd})

	ctx.notify(LevelError, "scope", bad)

	out := buf.String()

	if !strings.Contains(out, `"error"`) {
		t.Fatalf("expected fallback json error, got: %s", out)
	}
}

type failWriter struct{}

func (f failWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

func TestNotifyJSON_WriteErrorFallback(t *testing.T) {
	ctx := &context{
		w: &failWriter{},
		opt: Options{
			JSON: true,
		},
		color: NewColorizer(true),
	}

	ctx.notifyJSON(LevelInfo, "scope", "msg")
}
