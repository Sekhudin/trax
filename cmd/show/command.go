package show

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
		Use:     "show",
		Aliases: []string{"s"},
		Short:   "Inspect project information",
	},

	config: doc.Docs{
		Use:   "config",
		Short: "Show current Trax configuration",
		Long: doc.Paragraph(
			"Displays the resolved configuration used by Trax.",
		),
	},

	routes: doc.Docs{
		Use:   "routes",
		Short: "View registered routes",
		Long:  "Displays all registered routes.",
	},
}

func New(ctx *app.Context) *cobra.Command {
	cmd := doc.Apply(&docs.root, &cobra.Command{
		Args: cobra.NoArgs,
	})

	configCmd := NewConfigCmd(&docs.config, ctx)
	routesCmd := NewRoutesCmd(&docs.routes, ctx)

	cmd.AddCommand(configCmd, routesCmd)

	return cmd
}
