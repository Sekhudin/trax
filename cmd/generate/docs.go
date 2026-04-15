package generate

import "trax/internal/docs"

var (
	gDocs = docs.Docs{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate Typescript artifacts",
	}

	gRoutesDocs = docs.Docs{
		Use:   "routes",
		Short: "Generate type-safe route helpers",
		Long: docs.Paragraph(
			"Trax generates type-safe route helpers from your project structure.",
			"It allows you to navigate and construct routes programmatically with full type safety.",
		),
	}
)
