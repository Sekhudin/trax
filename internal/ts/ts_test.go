package ts

import "testing"

func TestGenerator_Success(t *testing.T) {
	t.Run("join_lines_newline", func(t *testing.T) {
		out := Line("a", "b", "c")
		if out != "a\nb\nc" {
			t.Fatal("newline_join_failed")
		}
	})

	t.Run("format_type_literal", func(t *testing.T) {
		if ToTypeLiteral("X") != "`X`" {
			t.Fatal("backtick_format_failed")
		}
	})

	t.Run("organize_unique_imports", func(t *testing.T) {
		mods := []string{"fs", "path", "fs"}
		if OrganizeImport(mods) != "fs, path" {
			t.Fatal("deduplicate_import_failed")
		}
	})
}

func TestUnion_Success(t *testing.T) {
	t.Run("single_union_value", func(t *testing.T) {
		if OrganizeStringUnion([]string{"A"}) != `"A"` {
			t.Fatal("single_format_failed")
		}
	})

	t.Run("multiple_unique_union", func(t *testing.T) {
		ss := []string{"A", "B", "A", "C"}
		expect := `"A" | "B" | "C"`
		if OrganizeStringUnion(ss) != expect {
			t.Fatal("union_deduplicate_failed")
		}
	})
}

func TestGenerator_Fallback(t *testing.T) {
	t.Run("empty_import_slice", func(t *testing.T) {
		if OrganizeImport([]string{}) != "" {
			t.Fatal("should_return_empty")
		}
	})

	t.Run("line_single_arg", func(t *testing.T) {
		if Line("only") != "only" {
			t.Fatal("single_line_failed")
		}
	})
}
