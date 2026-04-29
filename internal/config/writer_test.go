package config

import (
	"errors"
	"testing"
)

func TestWriter_Write_Success(t *testing.T) {
	t.Run("use_safe_write", func(t *testing.T) {
		called := false
		safe := func(f string) error {
			called = true
			if f != "a.toml" {
				t.Fatal("wrong_file")
			}
			return nil
		}
		force := func(string) error {
			t.Fatal("should_not_call_force")
			return nil
		}

		w := NewWriter("a.toml", false, safe, force)
		if err := w.Write(); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("safe_not_called")
		}
	})

	t.Run("use_force_write", func(t *testing.T) {
		called := false
		safe := func(string) error {
			t.Fatal("should_not_call_safe")
			return nil
		}
		force := func(f string) error {
			called = true
			if f != "b.toml" {
				t.Fatal("wrong_file")
			}
			return nil
		}

		w := NewWriter("b.toml", true, safe, force)
		if err := w.Write(); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("force_not_called")
		}
	})
}

func TestWriter_Write_Error(t *testing.T) {
	t.Run("propagate_write_error", func(t *testing.T) {
		expected := errors.New("boom")
		force := func(string) error { return expected }

		w := NewWriter("x.toml", true, nil, force)
		if err := w.Write(); !errors.Is(err, expected) {
			t.Fatal("error_not_propagated")
		}
	})
}

func TestWriter_File_Success(t *testing.T) {
	t.Run("return_correct_path", func(t *testing.T) {
		w := NewWriter("z.toml", false, nil, nil)
		if w.File() != "z.toml" {
			t.Fatal("wrong_file_path")
		}
	})
}
