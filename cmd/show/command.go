package show

import (
	"github.com/sekhudin/trax/internal/docs"

	"github.com/spf13/cobra"
)

type showdocs struct {
	root   docs.Docs
	config docs.Docs
	routes docs.Docs
}

var (
	doc = showdocs{
		root: docs.Docs{
			Use:     "show",
			Aliases: []string{"s"},
			Short:   "Inspect project information",
		},

		config: docs.Docs{
			Use:   "config",
			Short: "Show current Trax configuration",
			Long: docs.Paragraph(
				"Displays the resolved configuration used by Trax.",
			),
		},

		routes: docs.Docs{
			Use:   "routes",
			Short: "View registered routes",
			Long:  "Displays all registered routes.",
		},
	}

	Command = docs.ApplyDocs(&doc.root, &cobra.Command{})
)

func init() {
	Command.AddCommand(scCommand, srCommand)
}
