package doc

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLine(t *testing.T) {
	t.Run("multiple lines joined with newline", func(t *testing.T) {
		out := Line("a", "b", "c")
		expected := "a\nb\nc"

		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("single line returns as is", func(t *testing.T) {
		out := Line("only")
		if out != "only" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	t.Run("no lines returns empty string", func(t *testing.T) {
		out := Line()
		if out != "" {
			t.Fatalf("expected empty string, got %q", out)
		}
	})
}

func TestParagraph(t *testing.T) {
	t.Run("multiple lines joined with double newline", func(t *testing.T) {
		out := Paragraph("a", "b", "c")
		expected := "a\n\nb\n\nc"

		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("single line returns as is", func(t *testing.T) {
		out := Paragraph("only")
		if out != "only" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	t.Run("no lines returns empty string", func(t *testing.T) {
		out := Paragraph()
		if out != "" {
			t.Fatalf("expected empty string, got %q", out)
		}
	})
}

func TestApplyGroup(t *testing.T) {
	g := Group{
		ID:    "gen",
		Title: "Generate Commands",
	}

	cg := ApplyGroup(g)

	if cg.ID != "gen" {
		t.Fatalf("expected ID gen, got %s", cg.ID)
	}
	if cg.Title != "Generate Commands" {
		t.Fatalf("expected Title Generate Commands, got %s", cg.Title)
	}
}

func TestApply(t *testing.T) {
	t.Run("should apply all fields when provided", func(t *testing.T) {
		cmd := &cobra.Command{}

		d := &Docs{
			GroupID: "gen",
			Use:     "generate",
			Aliases: []string{"g"},
			Short:   "short desc",
			Long:    "long desc",
			Example: "example usage",
			Version: "1.0.0",
		}

		result := Apply(d, cmd)

		if result != cmd {
			t.Fatalf("expected same command pointer returned")
		}

		if cmd.GroupID != "gen" ||
			cmd.Use != "generate" ||
			cmd.Short != "short desc" ||
			cmd.Long != "long desc" ||
			cmd.Example != "example usage" ||
			cmd.Version != "1.0.0" {
			t.Fatalf("fields not properly applied")
		}

		if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "g" {
			t.Fatalf("aliases not applied")
		}
	})

	t.Run("should not overwrite existing fields when docs fields are empty", func(t *testing.T) {
		cmd := &cobra.Command{
			GroupID: "existing",
			Use:     "existing-use",
			Aliases: []string{"x"},
			Short:   "existing short",
			Long:    "existing long",
			Example: "existing example",
			Version: "0.9.0",
		}

		d := &Docs{}

		Apply(d, cmd)

		if cmd.GroupID != "existing" ||
			cmd.Use != "existing-use" ||
			cmd.Short != "existing short" ||
			cmd.Long != "existing long" ||
			cmd.Example != "existing example" ||
			cmd.Version != "0.9.0" {
			t.Fatalf("existing fields should not be overwritten")
		}

		if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "x" {
			t.Fatalf("aliases should not be overwritten")
		}
	})

	t.Run("should apply only non empty fields", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "old",
		}

		d := &Docs{
			Use:   "new",
			Short: "short",
		}

		Apply(d, cmd)

		if cmd.Use != "new" {
			t.Fatalf("Use should be updated")
		}
		if cmd.Short != "short" {
			t.Fatalf("Short should be set")
		}
	})
}
