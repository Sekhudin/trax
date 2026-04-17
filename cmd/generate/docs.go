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
			"Generate configuration file for your project",
			docs.Line(
				"Supported output formats:",
				"  - JSON (trax.json)",
				"  - YAML (trax.yaml)",
				"  - TOML (trax.toml)",
			),
		),
	}

	gRoutesDocs = docs.Docs{
		Use:   "routes",
		Short: "Generate type-safe route helpers",
		Long: docs.Paragraph(
			"Generates type-safe route helpers from your project structure or given routes declaration",
		),
	}
)
