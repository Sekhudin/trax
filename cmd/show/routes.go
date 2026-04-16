package show

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var sRoutesCmd = docs.ApplyDocs(sRoutesDocs, &cobra.Command{
	RunE: sRoutesRunE,
})

func init() {
}

func sRoutesRunE(cmd *cobra.Command, args []string) error {
	return nil
}
