package cmd

import (
	"trax/internal/docs"
)

var rootDocs = docs.Docs{
	Use:     "trax",
	Version: "0.0.1",
	Short:   "Powering TypeScript project workflows",
	Long: docs.Paragraph(
		"Trax is a CLI tool for automating TypeScript project workflows.",
		"It helps you generate and manage project structures with consistent conventions.",
	),
}
