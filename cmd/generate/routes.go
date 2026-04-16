package generate

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var gRoutesCmd = docs.ApplyDocs(gRoutesDocs, &cobra.Command{
	RunE: gRoutesRunE,
})

func init() {
}

func gRoutesRunE(cmd *cobra.Command, args []string) error {
	return nil
}
