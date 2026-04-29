package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestCoreError_Success(t *testing.T) {
	t.Run("format_only_code", func(t *testing.T) {
		e := &CoreError{Code: ErrInternal}
		if e.Error() != "[INTERNAL_ERROR]" {
			t.Fatal("wrong_code_format")
		}
	})

	t.Run("format_with_scope", func(t *testing.T) {
		e := &CoreError{Code: ErrIO, Scope: "file"}
		if !strings.Contains(e.Error(), "(file)") {
			t.Fatal("missing_scope_brackets")
		}
	})

	t.Run("format_full_combination", func(t *testing.T) {
		root := errors.New("boom")
		e := &CoreError{
			Code:    ErrExecution,
			Scope:   "cmd",
			Message: "failed",
			Err:     root,
		}
		out := e.Error()
		if !strings.Contains(out, "[EXECUTION_FAILED]") || !strings.Contains(out, "boom") {
			t.Fatal("incomplete_error_string")
		}
	})

	t.Run("unwrap_inner_error", func(t *testing.T) {
		root := errors.New("inner")
		e := &CoreError{Err: root}
		if !errors.Is(e, root) {
			t.Fatal("unwrap_failed")
		}
	})
}

func TestFactories_Success(t *testing.T) {
	root := errors.New("x")

	tests := []struct {
		name string
		err  error
		code ErrorCode
	}{
		{"config_load", NewConfigLoadError("cfg", "msg", root), ErrConfigLoad},
		{"config_missing", NewConfigNotFoundError("cfg", "msg"), ErrConfigNotFound},
		{"flag_error", NewFlagReadError("flag", root), ErrFlagRead},
		{"val_error", NewValidationError("v", "msg"), ErrValidation},
		{"io_error", NewIOError("io", "msg", root), ErrIO},
		{"tpl_error", NewTemplateError("tpl", "msg", root), ErrTemplate},
		{"rt_error", NewRuntimeError("rt", "msg", root), ErrRuntime},
		{"dep_error", NewDependencyError("dep", "msg", root), ErrDependency},
		{"int_error", NewInternalError("int", "msg", root), ErrInternal},
		{"cfg_invalid", NewInvalidConfigError("cfg", "msg"), ErrInvalidConfig},
		{"exec_error", NewExecutionError("exec", "msg", root), ErrExecution},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce, ok := tt.err.(*CoreError)
			if !ok || ce.Code != tt.code {
				t.Fatalf("invalid_code_assignment: %s", tt.code)
			}
		})
	}
}

func TestFactories_Fallback(t *testing.T) {
	t.Run("default_flag_message", func(t *testing.T) {
		err := NewFlagReadError("myflag", nil).(*CoreError)
		if err.Message != "failed to read flag value" {
			t.Fatal("wrong_default_msg")
		}
	})
}
