package doc

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLine_Success(t *testing.T) {
	t.Run("join_multiple_lines", func(t *testing.T) {
		out := Line("a", "b", "c")
		if out != "a\nb\nc" {
			t.Fatalf("wrong_line_output")
		}
	})

	t.Run("return_single_line", func(t *testing.T) {
		if Line("only") != "only" {
			t.Fatal("wrong_single_output")
		}
	})
}

func TestLine_Fallback(t *testing.T) {
	t.Run("return_empty_string", func(t *testing.T) {
		if Line() != "" {
			t.Fatal("expected_empty_string")
		}
	})
}

func TestParagraph_Success(t *testing.T) {
	t.Run("join_with_double_newline", func(t *testing.T) {
		out := Paragraph("a", "b")
		if out != "a\n\nb" {
			t.Fatal("wrong_paragraph_output")
		}
	})
}

func TestApplyGroup_Success(t *testing.T) {
	t.Run("map_group_fields", func(t *testing.T) {
		g := Group{ID: "gen", Title: "Generate"}
		cg := ApplyGroup(g)
		if cg.ID != "gen" || cg.Title != "Generate" {
			t.Fatal("fields_not_mapped")
		}
	})
}

func TestApply_Success(t *testing.T) {
	t.Run("apply_complete_docs", func(t *testing.T) {
		cmd := &cobra.Command{}
		d := &Docs{
			GroupID: "gen",
			Use:     "generate",
			Aliases: []string{"g", "gen"},
			Short:   "short desc",
			Long:    "long description",
			Example: "example usage",
			Version: "1.0.0",
		}

		Apply(d, cmd)

		if cmd.GroupID != d.GroupID || cmd.Use != d.Use || cmd.Version != d.Version {
			t.Fatal("failed_basic_fields")
		}
		if cmd.Short != d.Short || cmd.Long != d.Long || cmd.Example != d.Example {
			t.Fatal("failed_desc_fields")
		}
		if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "g" {
			t.Fatal("failed_aliases_mapping")
		}
	})

	t.Run("partial_update_fields", func(t *testing.T) {
		cmd := &cobra.Command{Use: "original", Short: "old"}
		d := &Docs{Use: "updated"}

		Apply(d, cmd)

		if cmd.Use != "updated" {
			t.Fatal("failed_use_update")
		}
		if cmd.Short != "old" {
			t.Fatal("overwritten_existing_short")
		}
	})
}

func TestApply_Fallback(t *testing.T) {
	t.Run("ignore_empty_docs", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:     "stay",
			Short:   "dont_change",
			Aliases: []string{"s"},
		}

		Apply(&Docs{}, cmd)

		if cmd.Use != "stay" || cmd.Short != "dont_change" || len(cmd.Aliases) != 1 {
			t.Fatal("data_was_overwritten")
		}
	})
}
