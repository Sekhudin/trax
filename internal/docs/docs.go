package docs

import (
	"strings"

	"github.com/spf13/cobra"
)

type Group struct {
	ID    string
	Title string
}

type Docs struct {
	GroupID string
	Use     string
	Aliases []string
	Version string
	Short   string
	Long    string
	Example string
}

func Line(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Paragraph(lines ...string) string {
	return strings.Join(lines, "\n\n")
}

func ApplyGroup(g Group) *cobra.Group {
	return &cobra.Group{
		ID:    g.ID,
		Title: g.Title,
	}
}

func ApplyDocs(d *Docs, cmd *cobra.Command) *cobra.Command {
	if d.GroupID != "" {
		cmd.GroupID = d.GroupID
	}

	if d.Use != "" {
		cmd.Use = d.Use
	}

	if len(d.Aliases) > 0 {
		cmd.Aliases = d.Aliases
	}

	if d.Short != "" {
		cmd.Short = d.Short
	}

	if d.Long != "" {
		cmd.Long = d.Long
	}

	if d.Example != "" {
		cmd.Example = d.Example
	}

	if d.Version != "" {
		cmd.Version = d.Version
	}

	return cmd
}
