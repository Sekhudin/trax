package docs

import (
	"strings"

	"github.com/spf13/cobra"
)

type Docs struct {
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

func ApplyDocs(d Docs, cmd *cobra.Command) *cobra.Command {
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
