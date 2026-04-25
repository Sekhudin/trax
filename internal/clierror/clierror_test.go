package clierror

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
	"github.com/sekhudin/trax/internal/output"
)

func newHandler(buf *bytes.Buffer) *Handler {
	ctx := output.New(buf, output.Options{NoColor: true})
	return New(ctx)
}

func TestHandler_Print_NilError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	h.Print(nil)

	if buf.Len() != 0 {
		t.Fatal("expected no output")
	}
}

func TestHandler_Print_NonCoreError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	h.Print(errors.New("boom"))

	out := buf.String()

	if !strings.Contains(out, "(runtime)") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "boom") {
		t.Fatal(out)
	}
}

func TestHandler_Print_CoreError_WithScope(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	err := &appErr.CoreError{
		Code:    appErr.ErrRuntime,
		Scope:   "exec",
		Message: "failed",
	}

	h.Print(err)

	out := buf.String()

	if !strings.Contains(out, "(exec)") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "failed") {
		t.Fatal(out)
	}
}

func TestHandler_Print_CoreError_EmptyScopeDefaultsToCore(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	err := &appErr.CoreError{
		Code:    appErr.ErrRuntime,
		Message: "oops",
	}

	h.Print(err)

	out := buf.String()

	if !strings.Contains(out, "(core)") {
		t.Fatal(out)
	}
}

func TestHandler_Print_CoreError_WithCause(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := newHandler(buf)

	root := errors.New("root cause")
	err := &appErr.CoreError{
		Code:    appErr.ErrRuntime,
		Scope:   "exec",
		Message: "failed",
		Err:     root,
	}

	h.Print(err)

	out := buf.String()

	if !strings.Contains(out, "root cause") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "(cause)") {
		t.Fatal(out)
	}
}

func TestHandler_ExitCode_Mapping(t *testing.T) {
	h := newHandler(bytes.NewBuffer(nil))

	tests := []struct {
		err  error
		want int
	}{
		{&appErr.CoreError{Code: appErr.ErrValidation}, 2},
		{&appErr.CoreError{Code: appErr.ErrConfigNotFound}, 3},
		{&appErr.CoreError{Code: appErr.ErrConfigLoad}, 4},
		{&appErr.CoreError{Code: appErr.ErrRuntime}, 1},
		{errors.New("x"), 1},
	}

	for _, tt := range tests {
		if h.ExitCode(tt.err) != tt.want {
			t.Fatalf("expected %d got %d", tt.want, h.ExitCode(tt.err))
		}
	}
}
