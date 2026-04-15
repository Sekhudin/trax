package errors

import (
	"errors"
	"testing"
)

func TestCoreError_ErrorFormat(t *testing.T) {
	err := &CoreError{
		Code:    ErrConfigLoad,
		Scope:   "config",
		Message: "failed to load",
		Err:     errors.New("disk error"),
	}

	got := err.Error()
	want := "[CONFIG_LOAD_FAILED] (config) failed to load: disk error"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCoreError_ErrorFormat_WithoutScope(t *testing.T) {
	err := &CoreError{
		Code:    ErrValidation,
		Message: "invalid input",
	}

	got := err.Error()
	want := "[VALIDATION_FAILED] invalid input"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCoreError_ErrorFormat_WithoutMessage(t *testing.T) {
	err := &CoreError{
		Code:  ErrIO,
		Scope: "file",
		Err:   errors.New("read failed"),
	}

	got := err.Error()
	want := "[IO_OPERATION_FAILED] (file): read failed"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCoreError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")

	err := &CoreError{
		Code: ErrRuntime,
		Err:  cause,
	}

	if !errors.Is(err, cause) {
		t.Fatalf("expected unwrap to match cause error")
	}
}

func TestNewConfigLoadError(t *testing.T) {
	cause := errors.New("disk error")

	err, ok := NewConfigLoadError("config", "failed", cause).(*CoreError)
	if !ok {
		t.Fatalf("expected *CoreError")
	}

	if err.Code != ErrConfigLoad {
		t.Fatalf("wrong code")
	}

	if err.Scope != "config" {
		t.Fatalf("wrong scope")
	}

	if err.Err != cause {
		t.Fatalf("wrong wrapped error")
	}
}

func TestNewConfigNotFoundError(t *testing.T) {
	err, ok := NewConfigNotFoundError("config", "not found").(*CoreError)
	if !ok {
		t.Fatalf("expected *CoreError")
	}

	if err.Code != ErrConfigNotFound {
		t.Fatalf("wrong code")
	}
}

func TestNewFlagReadError(t *testing.T) {
	cause := errors.New("invalid flag")

	err, ok := NewFlagReadError("debug", cause).(*CoreError)
	if !ok {
		t.Fatalf("expected *CoreError")
	}

	if err.Code != ErrFlagRead {
		t.Fatalf("wrong code")
	}

	if err.Scope != "debug" {
		t.Fatalf("wrong scope")
	}
}

func TestAllErrorBuilders_AreCoreError(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"validation", func() error { return NewValidationError("field", "invalid") }},
		{"io", func() error { return NewIOError("file", "failed", errors.New("x")) }},
		{"template", func() error { return NewTemplateError("tpl", "failed", errors.New("x")) }},
		{"runtime", func() error { return NewRuntimeError("exec", "failed", errors.New("x")) }},
		{"dependency", func() error { return NewDependencyError("pkg", "failed", errors.New("x")) }},
		{"internal", func() error { return NewInternalError("system", "failed", errors.New("x")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, ok := tt.fn().(*CoreError)
			if !ok {
				t.Fatalf("expected *CoreError")
			}

			if err.Code == "" {
				t.Fatalf("expected error code")
			}
		})
	}
}
