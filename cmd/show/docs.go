package show

import "trax/internal/docs"

var (
	Docs = docs.Docs{
		Use:     "show",
		Aliases: []string{"s"},
		Short:   "Inspect project information",
	}

	RoutesDocs = docs.Docs{
		Use:   "routes",
		Short: "View registered routes",
		Long:  "Displays all registered routes.",
	}

	ConfigDocs = docs.Docs{
		Use:   "config",
		Short: "Show current Trax configuration",
		Long: docs.Paragraph(
			"Displays the resolved configuration used by Trax.",
			"This includes CLI flags, config file values, and environment overrides.",
		),
	}
)
