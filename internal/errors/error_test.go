package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestCoreError_ErrorFormatting(t *testing.T) {
	t.Run("only code", func(t *testing.T) {
		e := &CoreError{Code: ErrInternal}
		if e.Error() != "[INTERNAL_ERROR]" {
			t.Fatal(e.Error())
		}
	})

	t.Run("with scope", func(t *testing.T) {
		e := &CoreError{Code: ErrIO, Scope: "file"}
		if !strings.Contains(e.Error(), "(file)") {
			t.Fatal(e.Error())
		}
	})

	t.Run("with message", func(t *testing.T) {
		e := &CoreError{Code: ErrValidation, Message: "bad input"}
		if !strings.Contains(e.Error(), "bad input") {
			t.Fatal(e.Error())
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		root := errors.New("root cause")
		e := &CoreError{Code: ErrRuntime, Err: root}
		if !strings.Contains(e.Error(), "root cause") {
			t.Fatal(e.Error())
		}
	})

	t.Run("full combination", func(t *testing.T) {
		root := errors.New("boom")
		e := &CoreError{
			Code:    ErrExecution,
			Scope:   "cmd",
			Message: "failed",
			Err:     root,
		}

		out := e.Error()

		if !strings.Contains(out, "[EXECUTION_FAILED]") ||
			!strings.Contains(out, "(cmd)") ||
			!strings.Contains(out, "failed") ||
			!strings.Contains(out, "boom") {
			t.Fatal(out)
		}
	})
}

func TestCoreError_Unwrap(t *testing.T) {
	root := errors.New("inner")
	e := &CoreError{Err: root}

	if !errors.Is(e, root) {
		t.Fatal("unwrap failed")
	}
}

func TestFactories_SetCorrectFields(t *testing.T) {
	root := errors.New("x")

	tests := []struct {
		name string
		err  error
		code ErrorCode
	}{
		{"config load", NewConfigLoadError("cfg", "msg", root), ErrConfigLoad},
		{"config not found", NewConfigNotFoundError("cfg", "msg"), ErrConfigNotFound},
		{"flag read", NewFlagReadError("flag", root), ErrFlagRead},
		{"validation", NewValidationError("v", "msg"), ErrValidation},
		{"io", NewIOError("io", "msg", root), ErrIO},
		{"template", NewTemplateError("tpl", "msg", root), ErrTemplate},
		{"runtime", NewRuntimeError("rt", "msg", root), ErrRuntime},
		{"dependency", NewDependencyError("dep", "msg", root), ErrDependency},
		{"internal", NewInternalError("int", "msg", root), ErrInternal},
		{"invalid config", NewInvalidConfigError("cfg", "msg"), ErrInvalidConfig},
		{"execution", NewExecutionError("exec", "msg", root), ErrExecution},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce, ok := tt.err.(*CoreError)
			if !ok {
				t.Fatal("not CoreError")
			}

			if ce.Code != tt.code {
				t.Fatalf("expected %s got %s", tt.code, ce.Code)
			}
		})
	}
}

func TestNewFlagReadError_DefaultMessage(t *testing.T) {
	root := errors.New("boom")
	err := NewFlagReadError("myflag", root).(*CoreError)

	if err.Message != "failed to read flag value" {
		t.Fatal(err.Message)
	}
}
