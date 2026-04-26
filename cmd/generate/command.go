package generate

import (
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"

	"github.com/spf13/cobra"
)

type documentation struct {
	root   doc.Docs
	config doc.Docs
	routes doc.Docs
}

var docs = documentation{
	root: doc.Docs{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate artifacts",
	},

	config: doc.Docs{
		Use:   "config",
		Short: "Generate configuration file",
		Long: doc.Paragraph(
			"Generate configuration file for your project",
			doc.Line(
				"Supported output formats:",
				"  - JSON (trax.json)",
				"  - YAML (trax.yaml)",
				"  - TOML (trax.toml)",
			),
		),
	},

	routes: doc.Docs{
		Use:   "routes",
		Short: "Generate type-safe route helpers",
		Long: doc.Paragraph(
			"Generates type-safe route helpers from your project structure or given routes declaration",
		),
	},
}

func New(ctx *app.Context) *cobra.Command {
	cmd := doc.Apply(&docs.root, &cobra.Command{
		Args: cobra.ExactArgs(1),
	})

	configCmd := NewConfigCmd(&docs.config, ctx)
	routesCmd := NewRoutesCmd(&docs.routes, ctx)

	cmd.AddCommand(configCmd, routesCmd)

	return cmd
}
