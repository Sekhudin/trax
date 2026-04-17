package show

import "trax/internal/docs"

var (
	sDocs = docs.Docs{
		Use:     "show",
		Aliases: []string{"s"},
		Short:   "Inspect project information",
	}

	sRoutesDocs = docs.Docs{
		Use:   "routes",
		Short: "View registered routes",
		Long:  "Displays all registered routes.",
	}

	sConfigDocs = docs.Docs{
		Use:   "config",
		Short: "Show current Trax configuration",
		Long: docs.Paragraph(
			"Displays the resolved configuration used by Trax.",
		),
	}
)
