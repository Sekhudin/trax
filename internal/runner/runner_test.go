package runner

import (
	"bytes"
	"errors"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
	"github.com/sekhudin/trax/internal/testutil/errormock"
)

func TestRunner_Success(t *testing.T) {
	t.Run("execute_echo_command", func(t *testing.T) {
		buf := &bytes.Buffer{}
		r := NewRunner(buf, buf)

		cmd := map[string]any{
			"exec": "echo",
			"args": []any{"trax", 1, true},
		}

		if err := r.Run(cmd); err != nil {
			t.Fatalf("unexpected_error: %v", err)
		}
	})

	t.Run("args_string_slice", func(t *testing.T) {
		r := &CommandRunner{}
		_, args, err := r.parseCommand(map[string]any{
			"exec": "ls",
			"args": []string{"-a", "-l"},
		})
		if err != nil || len(args) != 2 {
			t.Fatal("parse_slice_failed")
		}
	})
}

func TestRunner_Error(t *testing.T) {
	r := &CommandRunner{}

	t.Run("nil_command_config", func(t *testing.T) {
		err := r.Run(nil)
		if err == nil || !errors.Is(err, errormock.With(appErr.ErrInvalidConfig)) {
			t.Fatal("should_invalid_config")
		}
	})

	t.Run("missing_exec_field", func(t *testing.T) {
		err := r.Run(map[string]any{"args": []string{"hi"}})
		if err == nil || !errors.Is(err, errormock.With(appErr.ErrInvalidConfig)) {
			t.Fatal("should_require_exec")
		}
	})

	t.Run("command_not_found", func(t *testing.T) {
		r := NewRunner(&bytes.Buffer{}, &bytes.Buffer{})
		err := r.Run(map[string]any{"exec": "invalid-cmd-123"})
		if err == nil || !errors.Is(err, errormock.With(appErr.ErrExecution)) {
			t.Fatal("should_execution_error")
		}
	})
}

func TestRunner_Fallback(t *testing.T) {
	t.Run("interface_implementation_check", func(t *testing.T) {
		var _ Runner = (*CommandRunner)(nil)
	})

	t.Run("writer_wiring_check", func(t *testing.T) {
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		r := NewRunner(stdout, stderr).(*CommandRunner)
		if r.Stdout != stdout || r.Stderr != stderr {
			t.Fatal("wiring_failed")
		}
	})

	t.Run("mixed_type_args", func(t *testing.T) {
		r := &CommandRunner{}
		_, args, _ := r.parseCommand(map[string]any{
			"exec": "echo",
			"args": []any{3.14, false},
		})
		if args[0] != "3.14" || args[1] != "false" {
			t.Fatal("string_conversion_failed")
		}
	})
}
