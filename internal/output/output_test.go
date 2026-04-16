package output_test

import (
	"bytes"
	"strings"
	"testing"

	"trax/internal/output"

	"github.com/spf13/viper"
)

func reset(t *testing.T) {
	t.Helper()
	viper.Reset()
}

func TestContext_Success_TextMode(t *testing.T) {
	reset(t)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Success("config", "created")

	got := buf.String()

	if !strings.Contains(got, "✔ [config] created") {
		t.Errorf("unexpected output: %s", got)
	}
}

func TestContext_Info_TextMode(t *testing.T) {
	reset(t)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Info("auth", "login success")

	got := buf.String()

	if !strings.Contains(got, "ℹ  [auth] login success") {
		t.Errorf("unexpected output: %s", got)
	}
}

func TestContext_Warn_TextMode(t *testing.T) {
	reset(t)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Warn("config", "missing file")

	got := buf.String()

	if !strings.Contains(got, "  [config] missing file") {
		t.Errorf("unexpected output: %s", got)
	}
}

func TestContext_Success_DebugMode_JSON(t *testing.T) {
	reset(t)
	viper.Set("debug", true)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Success("config", "created")

	got := buf.String()

	if !strings.Contains(got, `"level": "success"`) {
		t.Errorf("expected json output, got: %s", got)
	}

	if !strings.Contains(got, `"scope": "config"`) {
		t.Errorf("missing scope in json: %s", got)
	}

	if !strings.Contains(got, `"message": "created"`) {
		t.Errorf("missing message in json: %s", got)
	}
}

func TestContext_Info_DebugMode_JSON(t *testing.T) {
	reset(t)
	viper.Set("debug", true)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Info("auth", "ok")

	got := buf.String()

	if !strings.Contains(got, `"level": "info"`) {
		t.Errorf("expected info json, got: %s", got)
	}
}

func TestContext_Warn_DebugMode_JSON(t *testing.T) {
	reset(t)
	viper.Set("debug", true)

	var buf bytes.Buffer
	out := output.New(&buf)

	out.Warn("config", "missing")

	got := buf.String()

	if !strings.Contains(got, `"level": "warn"`) {
		t.Errorf("expected warn json, got: %s", got)
	}
}
