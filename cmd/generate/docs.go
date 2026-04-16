package generate

import "trax/internal/docs"

var (
	gDocs = docs.Docs{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate Typescript artifacts",
	}

	gConfigDocs = docs.Docs{
		Use:   "config",
		Short: "Generate configuration file",
		Long: docs.Paragraph(
			"Trax generates a configuration file for your project setup and CLI behavior.",
			docs.Line(
				"Supported output formats:",
				"- JSON (trax.json)",
				"- YAML (trax.yaml)",
				"- TOML (trax.toml)",
			),
		),
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
