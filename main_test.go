package main

import (
	"io"
	"os"
	"testing"
)

func TestMain_Success(t *testing.T) {
	t.Run("main_normal_flow", func(t *testing.T) {
		main()
	})
}

func TestRecover_Logic(t *testing.T) {
	t.Run("handle_panic_recovery", func(t *testing.T) {
		oldExit := Exit
		exitCalled := false
		Exit = func(code int) { exitCalled = true }
		defer func() { Exit = oldExit }()

		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		func() {
			defer Recover()
			panic("system failure")
		}()

		w.Close()
		os.Stderr = oldStderr
		out, _ := io.ReadAll(r)

		if !exitCalled {
			t.Fatal("exit not called")
		}

		expected := "[INTERNAL ERROR] system failure\n"
		if string(out) != expected {
			t.Errorf("got %q, want %q", string(out), expected)
		}
	})

	t.Run("no_panic_flow", func(t *testing.T) {
		oldExit := Exit
		exitCalled := false
		Exit = func(code int) { exitCalled = true }
		defer func() { Exit = oldExit }()

		func() {
			defer Recover()
		}()

		if exitCalled {
			t.Fatal("exit called unnecessarily")
		}
	})
}
