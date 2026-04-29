package clierror

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
	"github.com/sekhudin/trax/internal/output"
)

func newHandler(buf *bytes.Buffer) Handler {
	ctx := output.New(buf, output.Options{NoColor: true})
	return New(ctx)
}

func TestHandler_Print_Success(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	t.Run("print_core_error", func(t *testing.T) {
		buf.Reset()
		err := &appErr.CoreError{Scope: "exec", Message: "failed"}
		h.Print(err)
		if !strings.Contains(buf.String(), "(exec)") {
			t.Fatal("missing_scope_output")
		}
	})

	t.Run("print_with_cause", func(t *testing.T) {
		buf.Reset()
		err := &appErr.CoreError{Message: "fail", Err: errors.New("root")}
		h.Print(err)
		if !strings.Contains(buf.String(), "root") {
			t.Fatal("missing_cause_output")
		}
	})

	t.Run("print_runtime_error", func(t *testing.T) {
		buf.Reset()
		h.Print(errors.New("boom"))
		if !strings.Contains(buf.String(), "(runtime)") {
			t.Fatal("missing_runtime_label")
		}
	})
}

func TestHandler_Print_Fallback(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	t.Run("handle_nil_error", func(t *testing.T) {
		buf.Reset()
		h.Print(nil)
		if buf.Len() != 0 {
			t.Fatal("expected_no_output")
		}
	})

	t.Run("default_to_core", func(t *testing.T) {
		buf.Reset()
		err := &appErr.CoreError{Message: "oops"}
		h.Print(err)
		if !strings.Contains(buf.String(), "(core)") {
			t.Fatal("missing_core_default")
		}
	})
}

func TestHandler_ExitCode_Success(t *testing.T) {
	h := newHandler(bytes.NewBuffer(nil))

	t.Run("map_validation_code", func(t *testing.T) {
		code := h.ExitCode(&appErr.CoreError{Code: appErr.ErrValidation})
		if code != 2 {
			t.Fatal("wrong_exit_code")
		}
	})

	t.Run("map_config_codes", func(t *testing.T) {
		if h.ExitCode(&appErr.CoreError{Code: appErr.ErrConfigNotFound}) != 3 {
			t.Fatal("wrong_notfound_code")
		}
		if h.ExitCode(&appErr.CoreError{Code: appErr.ErrConfigLoad}) != 4 {
			t.Fatal("wrong_load_code")
		}
	})
}

func TestHandler_ExitCode_Fallback(t *testing.T) {
	h := newHandler(bytes.NewBuffer(nil))

	t.Run("default_exit_code", func(t *testing.T) {
		if h.ExitCode(errors.New("x")) != 1 {
			t.Fatal("should_be_1")
		}
		if h.ExitCode(&appErr.CoreError{Code: appErr.ErrRuntime}) != 1 {
			t.Fatal("should_be_1")
		}
	})
}
