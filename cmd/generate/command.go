package generate

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

type generatedocs struct {
	root   docs.Docs
	config docs.Docs
	routes docs.Docs
}

var (
	doc = generatedocs{
		root: docs.Docs{
			Use:     "generate",
			Aliases: []string{"g"},
			Short:   "Generate Typescript artifacts",
		},

		config: docs.Docs{
			Use:   "config",
			Short: "Generate configuration file",
			Long: docs.Paragraph(
				"Generate configuration file for your project",
				docs.Line(
					"Supported output formats:",
					"  - JSON (trax.json)",
					"  - YAML (trax.yaml)",
					"  - TOML (trax.toml)",
				),
			),
		},

		routes: docs.Docs{
			Use:   "routes",
			Short: "Generate type-safe route helpers",
			Long: docs.Paragraph(
				"Generates type-safe route helpers from your project structure or given routes declaration",
			),
		},
	}

	Command = docs.ApplyDocs(&doc.root, &cobra.Command{})
)

func init() {
	Command.AddCommand(gcCommand, grCommand)
}
