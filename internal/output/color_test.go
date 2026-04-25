package output

import (
	"strings"
	"testing"
)

func TestColorizer_NoColor(t *testing.T) {
	c := NewColorizer(true)

	t.Run("all methods return plain text", func(t *testing.T) {
		input := []any{"hello", " ", "world"}

		tests := []struct {
			name string
			fn   func(...any) string
		}{
			{"red", c.Red},
			{"yellow", c.Yellow},
			{"green", c.Green},
			{"blue", c.Blue},
			{"cyan", c.Cyan},
			{"gray", c.Gray},
			{"bold", c.Bold},
		}

		for _, tt := range tests {
			out := tt.fn(input...)
			if out != "hello world" {
				t.Fatalf("%s unexpected: %s", tt.name, out)
			}
		}
	})
}

func TestColorizer_WithColor(t *testing.T) {
	c := NewColorizer(false)

	t.Run("red contains ansi code", func(t *testing.T) {
		out := c.Red("fail")

		if !strings.Contains(out, "\x1b[31m") {
			t.Fatal(out)
		}
		if !strings.Contains(out, "fail") {
			t.Fatal(out)
		}
		if !strings.HasSuffix(out, "\x1b[0m") {
			t.Fatal(out)
		}
	})

	t.Run("bold uses code 1", func(t *testing.T) {
		out := c.Bold("x")

		if !strings.Contains(out, "\x1b[1m") {
			t.Fatal(out)
		}
	})

	t.Run("multiple values are sprinted", func(t *testing.T) {
		out := c.Green("a", 1, "b")

		if !strings.Contains(out, "a1b") {
			t.Fatal(out)
		}
	})
}
