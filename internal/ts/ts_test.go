package ts

import "testing"

func TestLine(t *testing.T) {
	out := Line("a", "b", "c")
	expected := "a\nb\nc"

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestTypeAndStringLiteral(t *testing.T) {
	t.Run("type literal", func(t *testing.T) {
		out := ToTypeLiteral("hello")
		if out != "`hello`" {
			t.Fatalf("unexpected: %s", out)
		}
	})

	t.Run("string union", func(t *testing.T) {
		out := ToStringUnion("hello")
		if out != `"hello"` {
			t.Fatalf("unexpected: %s", out)
		}
	})
}

func TestOrganizeImport(t *testing.T) {
	mods := []string{
		"fs",
		"path",
		"fs", // duplicate
		"http",
		"path", // duplicate
	}

	out := OrganizeImport(mods)

	// urutan tetap sesuai first occurrence
	expected := "fs, path, http"
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestOrganizeStringUnion(t *testing.T) {
	t.Run("single value", func(t *testing.T) {
		out := OrganizeStringUnion([]string{"A"})
		if out != `"A"` {
			t.Fatalf("unexpected: %s", out)
		}
	})

	t.Run("multiple with duplicates", func(t *testing.T) {
		ss := []string{"A", "B", "A", "C", "B"}

		out := OrganizeStringUnion(ss)

		// urutan sesuai first occurrence
		expected := `"A" | "B" | "C"`
		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})
}
