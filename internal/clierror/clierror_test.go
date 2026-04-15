package clierror

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

func reset(t *testing.T) {
	t.Helper()
	viper.Reset()
	viper.SetDefault("debug", false)
}

func TestPrintText_CoreError(t *testing.T) {
	reset(t)

	var buf bytes.Buffer

	err := &appErr.CoreError{
		Code:    appErr.ErrConfigLoad,
		Scope:   "config",
		Message: "failed to load",
		Err:     errors.New("disk error"),
	}

	printText(&buf, err)

	out := buf.String()

	if !strings.Contains(out, "[CONFIG_LOAD_FAILED]") {
		t.Fatal("missing code")
	}

	if !strings.Contains(out, "(config)") {
		t.Fatal("missing scope")
	}

	if !strings.Contains(out, "failed to load") {
		t.Fatal("missing message")
	}

	if !strings.Contains(out, "disk error") {
		t.Fatal("missing cause")
	}
}

func TestPrintText_DefaultError(t *testing.T) {
	reset(t)

	var buf bytes.Buffer

	err := errors.New("generic error")

	printText(&buf, err)

	out := buf.String()

	if !strings.Contains(out, "❌ generic error") {
		t.Fatal("expected fallback error output")
	}
}

func TestPrintJSON_CoreError(t *testing.T) {
	reset(t)

	var buf bytes.Buffer

	err := &appErr.CoreError{
		Code:    appErr.ErrValidation,
		Scope:   "field",
		Message: "invalid input",
		Err:     errors.New("bad format"),
	}

	printJSON(&buf, err)

	out := buf.String()

	if !strings.Contains(out, `"code": "VALIDATION_FAILED"`) {
		t.Fatal("missing code")
	}

	if !strings.Contains(out, `"scope": "field"`) {
		t.Fatal("missing scope")
	}

	if !strings.Contains(out, `"cause": "bad format"`) {
		t.Fatal("missing cause")
	}
}

func TestPrint_TextMode(t *testing.T) {
	reset(t)
	viper.Set("debug", false)

	var buf bytes.Buffer

	err := &appErr.CoreError{
		Code:    appErr.ErrIO,
		Scope:   "file",
		Message: "failed",
	}

	Print(&buf, err)

	out := buf.String()

	if !strings.Contains(out, "[IO_OPERATION_FAILED]") {
		t.Fatal("expected text output")
	}

	if strings.Contains(out, "{") {
		t.Fatal("should not be JSON output")
	}
}

func TestPrint_JSONMode(t *testing.T) {
	reset(t)
	viper.Set("debug", true)

	var buf bytes.Buffer

	err := &appErr.CoreError{
		Code:    appErr.ErrRuntime,
		Scope:   "exec",
		Message: "crash",
		Err:     errors.New("panic"),
	}

	Print(&buf, err)

	out := buf.String()

	if !strings.Contains(out, `"code"`) {
		t.Fatal("expected json output")
	}

	if !strings.Contains(out, "RUNTIME_EXECUTION_FAILED") {
		t.Fatal("missing code in json")
	}
}

func TestExitCode(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{
			"validation",
			&appErr.CoreError{Code: appErr.ErrValidation},
			2,
		},
		{
			"config_not_found",
			&appErr.CoreError{Code: appErr.ErrConfigNotFound},
			3,
		},
		{
			"config_load",
			&appErr.CoreError{Code: appErr.ErrConfigLoad},
			4,
		},
		{
			"default_core_error",
			&appErr.CoreError{Code: "UNKNOWN"},
			1,
		},
		{
			"non_core_error",
			errors.New("random"),
			1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExitCode(c.err)

			if got != c.want {
				t.Fatalf("expected %d, got %d", c.want, got)
			}
		})
	}
}
