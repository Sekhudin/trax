package clierror

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
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
		t.Fatal("writer not set")
	}
}

func TestPrintText_CoreError_Simple(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	err := &appErr.CoreError{
		Code:    appErr.ErrConfigLoad,
		Message: "failed",
	}

	ctx.PrintText(err)

	out := buf.String()
	if !strings.Contains(out, "CONFIG_LOAD_FAILED") {
		t.Fatal(out)
	}
}

func TestPrintText_CoreError_WithScopeAndCause(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	err := &appErr.CoreError{
		Code:    appErr.ErrConfigLoad,
		Scope:   "config",
		Message: "failed",
		Err:     errors.New("root cause"),
	}

	ctx.PrintText(err)

	out := buf.String()
	if !strings.Contains(out, "config") ||
		!strings.Contains(out, "root cause") {
		t.Fatal(out)
	}
}

func TestPrintText_DefaultError(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.PrintText(errors.New("boom"))

	if !strings.Contains(buf.String(), "boom") {
		t.Fatal(buf.String())
	}
}

func TestPrintJSON_CoreError(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	err := &appErr.CoreError{
		Code:    appErr.ErrConfigLoad,
		Scope:   "config",
		Message: "failed",
		Err:     errors.New("root cause"),
	}

	ctx.PrintJSON(err)

	out := buf.String()

	if !strings.Contains(out, `"code"`) ||
		!strings.Contains(out, `"scope"`) ||
		!strings.Contains(out, `"cause"`) {
		t.Fatal(out)
	}
}

func TestPrintJSON_DefaultError(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.PrintJSON(errors.New("boom"))

	if !strings.Contains(buf.String(), `"error": "boom"`) {
		t.Fatal(buf.String())
	}
}

func TestPrint_TextMode(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.Print(errors.New("boom"))

	if !strings.Contains(buf.String(), "boom") {
		t.Fatal(buf.String())
	}
}

func TestPrint_DebugMode(t *testing.T) {
	resetViper(t)
	viper.Set("debug", true)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.Print(errors.New("boom"))

	if !strings.Contains(buf.String(), `"error": "boom"`) {
		t.Fatal(buf.String())
	}
}

func TestPrint_NilError(t *testing.T) {
	resetViper(t)

	buf := new(bytes.Buffer)
	ctx := New(buf)

	ctx.Print(nil)

	if buf.Len() != 0 {
		t.Fatal("should not print anything")
	}
}

func TestExitCode(t *testing.T) {
	resetViper(t)

	ctx := New(new(bytes.Buffer))

	tests := []struct {
		err  error
		want int
	}{
		{&appErr.CoreError{Code: appErr.ErrValidation}, 2},
		{&appErr.CoreError{Code: appErr.ErrConfigNotFound}, 3},
		{&appErr.CoreError{Code: appErr.ErrConfigLoad}, 4},
		{&appErr.CoreError{Code: "OTHER"}, 1},
		{errors.New("x"), 1},
	}

	for _, tt := range tests {
		got := ctx.ExitCode(tt.err)
		if got != tt.want {
			t.Fatalf("want %d got %d", tt.want, got)
		}
	}
}
