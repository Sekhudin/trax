package generate

import (
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"

	"github.com/spf13/cobra"
)

type Docs struct {
	Root   doc.Docs
	Config doc.Docs
	Routes doc.Docs
}

type Dependencies struct {
	Docs         *Docs
	NewConfigCtx func(app.Context) ConfigCtx
	NewRoutesCtx func(app.Context) RoutesCtx
}

func New(ctx app.Context) *cobra.Command {
	return NewWithDependencies(ctx, DefaultDependencies())
}

func NewWithDependencies(ctx app.Context, d *Dependencies) *cobra.Command {
	cmd := doc.Apply(&d.Docs.Root, &cobra.Command{
		Args: cobra.NoArgs,
	})

	configCmd := NewConfigCmd(&d.Docs.Config, d.NewConfigCtx(ctx))
	routesCmd := NewRoutesCmd(&d.Docs.Routes, d.NewRoutesCtx(ctx))

	cmd.AddCommand(configCmd, routesCmd)

	return cmd
}

func DefaultDependencies() *Dependencies {
	return &Dependencies{
		NewConfigCtx: NewConfigCtx,
		NewRoutesCtx: NewRoutesCtx,
		Docs: &Docs{
			Root: doc.Docs{
				Use:     "generate",
				Aliases: []string{"g"},
				Short:   "Generate artifacts",
			},

			Config: doc.Docs{
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

			Routes: doc.Docs{
				Use:   "routes",
				Short: "Generate type-safe route helpers",
				Long: doc.Paragraph(
					"Generates type-safe route helpers from your project structure or given routes declaration",
				),
			},
		},
	}
}
