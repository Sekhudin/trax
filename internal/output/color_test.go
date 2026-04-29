package output

import (
	"strings"
	"testing"
)

func TestColorizer_Success(t *testing.T) {
	c := NewColorizer(false)

	t.Run("wrap_ansi_codes", func(t *testing.T) {
		out := c.Red("fail")
		if !strings.HasPrefix(out, "\x1b[31m") || !strings.Contains(out, "fail") || !strings.HasSuffix(out, "\x1b[0m") {
			t.Fatal("wrong_ansi_formatting")
		}
	})

	t.Run("apply_bold_code", func(t *testing.T) {
		if !strings.Contains(c.Bold("x"), "\x1b[1m") {
			t.Fatal("wrong_bold_code")
		}
	})

	t.Run("sprint_multiple_values", func(t *testing.T) {
		out := c.Green("a", 1, "b")
		if !strings.Contains(out, "a1b") {
			t.Fatal("failed_join_values")
		}
	})
}

func TestColorizer_Fallback(t *testing.T) {
	c := NewColorizer(true)

	t.Run("return_plain_text", func(t *testing.T) {
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
			t.Run(tt.name, func(t *testing.T) {
				if out := tt.fn(input...); out != "hello world" {
					t.Fatalf("expected_plain_got_%q", out)
				}
			})
		}
	})
}
