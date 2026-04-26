package runner

import (
	"bytes"
	"errors"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func TestParseCommand(t *testing.T) {
	r := &CommandRunner{}

	t.Run("nil command -> invalid config error", func(t *testing.T) {
		_, _, err := r.parseCommand(nil)
		if err == nil {
			t.Fatal("expected error")
		}

		var ce *appErr.CoreError
		if !errors.As(err, &ce) {
			t.Fatal("expected CoreError")
		}

		if ce.Code != appErr.ErrInvalidConfig {
			t.Fatalf("unexpected code: %s", ce.Code)
		}
	})

	t.Run("missing exec field", func(t *testing.T) {
		_, _, err := r.parseCommand(map[string]any{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid exec type", func(t *testing.T) {
		_, _, err := r.parseCommand(map[string]any{
			"exec": 123,
		})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty exec string", func(t *testing.T) {
		_, _, err := r.parseCommand(map[string]any{
			"exec": "",
		})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("args as []any", func(t *testing.T) {
		exe, args, err := r.parseCommand(map[string]any{
			"exec": "echo",
			"args": []any{"a", 1, true},
		})
		if err != nil {
			t.Fatal(err)
		}

		if exe != "echo" || len(args) != 3 {
			t.Fatalf("unexpected: %s %#v", exe, args)
		}
	})

	t.Run("args as []string", func(t *testing.T) {
		_, args, err := r.parseCommand(map[string]any{
			"exec": "echo",
			"args": []string{"a", "b"},
		})
		if err != nil {
			t.Fatal(err)
		}

		if args[0] != "a" || args[1] != "b" {
			t.Fatalf("unexpected args: %#v", args)
		}
	})

	t.Run("no args", func(t *testing.T) {
		exe, args, err := r.parseCommand(map[string]any{
			"exec": "echo",
		})
		if err != nil {
			t.Fatal(err)
		}

		if exe != "echo" || len(args) != 0 {
			t.Fatalf("unexpected: %s %#v", exe, args)
		}
	})
}

func TestRun_Success(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	r := NewRunner(stdout, stderr).(*CommandRunner)

	err := r.Run(map[string]any{
		"exec": "echo",
		"args": []string{"hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_Failure_CommandNotFound(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	r := NewRunner(stdout, stderr).(*CommandRunner)

	err := r.Run(map[string]any{
		"exec": "this-command-does-not-exist-xyz",
	})

	if err == nil {
		t.Fatal("expected error")
	}

	var ce *appErr.CoreError
	if !errors.As(err, &ce) {
		t.Fatal("expected CoreError")
	}

	if ce.Code != appErr.ErrExecution {
		t.Fatalf("expected execution error, got %s", ce.Code)
	}
}

func TestRun_ParseError_Propagation(t *testing.T) {
	r := &CommandRunner{}

	err := r.Run(nil)

	if err == nil {
		t.Fatal("expected error")
	}

	var ce *appErr.CoreError
	if !errors.As(err, &ce) {
		t.Fatal("expected CoreError from parseCommand")
	}

	if ce.Code != appErr.ErrInvalidConfig {
		t.Fatalf("expected invalid config, got %s", ce.Code)
	}
}

func TestInterface_Implementation(t *testing.T) {
	var _ Runner = &CommandRunner{}
}

func TestWriterIsWired(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	r := NewRunner(stdout, stderr).(*CommandRunner)

	if r.Stdout != stdout || r.Stderr != stderr {
		t.Fatal("writer not wired correctly")
	}
}

func TestParseCommand_ArgsMixedEdge(t *testing.T) {
	r := &CommandRunner{}

	_, args, err := r.parseCommand(map[string]any{
		"exec": "echo",
		"args": []any{1, "x", false, 3.14},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(args) != 4 {
		t.Fatalf("unexpected args: %#v", args)
	}
}
