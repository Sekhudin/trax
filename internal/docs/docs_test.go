package docs

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLine(t *testing.T) {
	got := Line("a", "b", "c")
	want := "a\nb\nc"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParagraph(t *testing.T) {
	got := Paragraph("a", "b", "c")
	want := "a\n\nb\n\nc"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApplyGroup(t *testing.T) {
	g := Group{
		ID:    "gen",
		Title: "Generate Commands",
	}

	out := ApplyGroup(g)

	if out.ID != "gen" {
		t.Fatal("expected group id")
	}

	if out.Title != "Generate Commands" {
		t.Fatal("expected group title")
	}
}

func TestApplyDocs_AllFields(t *testing.T) {
	cmd := &cobra.Command{}

	d := Docs{
		GroupID: "gen",
		Use:     "trax generate",
		Aliases: []string{"g"},
		Version: "1.0.0",
		Short:   "short desc",
		Long:    "long desc",
		Example: "example",
	}

	ApplyDocs(&d, cmd)

	if cmd.GroupID != "gen" {
		t.Fatal("group id not set")
	}

	if cmd.Use != "trax generate" {
		t.Fatal("use not set")
	}

	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "g" {
		t.Fatal("aliases not set")
	}

	if cmd.Version != "1.0.0" {
		t.Fatal("version not set")
	}

	if cmd.Short != "short desc" {
		t.Fatal("short not set")
	}

	if cmd.Long != "long desc" {
		t.Fatal("long not set")
	}

	if cmd.Example != "example" {
		t.Fatal("example not set")
	}
}

func TestApplyDocs_IgnoresEmptyFields(t *testing.T) {
	cmd := &cobra.Command{
		Use:     "original",
		Short:   "keep me",
		Version: "0.1.0",
	}

	d := Docs{
		Use: "", // should NOT override
	}

	ApplyDocs(&d, cmd)

	if cmd.Use != "original" {
		t.Fatal("empty use should not override")
	}

	if cmd.Short != "keep me" {
		t.Fatal("empty docs should not override existing")
	}

	if cmd.Version != "0.1.0" {
		t.Fatal("version should remain unchanged")
	}
}

func TestApplyDocs_PartialOverride(t *testing.T) {
	cmd := &cobra.Command{
		Short: "old short",
		Long:  "old long",
	}

	d := Docs{
		Short: "new short",
	}

	ApplyDocs(&d, cmd)

	if cmd.Short != "new short" {
		t.Fatal("short should be overridden")
	}

	if cmd.Long != "old long" {
		t.Fatal("long should remain unchanged")
	}
}
