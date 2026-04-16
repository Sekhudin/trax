package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func resetViper(t *testing.T) {
	t.Helper()
	viper.Reset()
}

func TestNew(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	if ctx.Writer != buf {
		t.Fatal("writer not set correctly")
	}
}

func TestNotifyText(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.NotifyText(Notification{"success", "a", "b"})
	ctx.NotifyText(Notification{"warn", "c", "d"})
	ctx.NotifyText(Notification{"info", "e", "f"})

	out := buf.String()

	if !strings.Contains(out, "✔ [a] b") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "  [c] d") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "ℹ  [e] f") {
		t.Fatal(out)
	}
}

func TestNotifyJSON(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.NotifyJSON(Notification{"success", "config", "ok"})

	out := buf.String()

	if !strings.Contains(out, `"level": "success"`) {
		t.Fatal(out)
	}
	if !strings.Contains(out, `"scope": "config"`) {
		t.Fatal(out)
	}
}

func TestNotifications_TextMode(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.Success("config", "ok")
	ctx.Info("config", "info")
	ctx.Warn("config", "warn")

	out := buf.String()

	if !strings.Contains(out, "✔ [config] ok") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "ℹ  [config] info") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "  [config] warn") {
		t.Fatal(out)
	}
}

func TestNotifications_DebugMode(t *testing.T) {
	resetViper(t)
	viper.Set("debug", true)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.Success("config", "ok")
	ctx.Info("config", "info")
	ctx.Warn("config", "warn")

	out := buf.String()

	if strings.Count(out, `"level"`) != 3 {
		t.Fatal(out)
	}
}

func TestAsJSON_OK(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	err := ctx.AsJSON(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), `"a": 1`) {
		t.Fatal(buf.String())
	}
}

func TestAsJSON_MarshalError(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	err := ctx.AsJSON(map[string]any{
		"bad": make(chan int),
	})

	if err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestAsFlat_NestedMap(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	data := map[string]any{
		"db": map[string]any{
			"host": "localhost",
		},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "db.host = localhost") {
		t.Fatal(buf.String())
	}
}

func TestAsFlat_SlicePrimitive(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	data := map[string]any{
		"nums": []any{1, 2},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	if !strings.Contains(out, "nums.0 = 1") ||
		!strings.Contains(out, "nums.1 = 2") {
		t.Fatal(out)
	}
}

func TestAsFlat_SliceMap(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	data := map[string]any{
		"items": []any{
			map[string]any{"a": 1},
		},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "items.0.a = 1") {
		t.Fatal(buf.String())
	}
}
