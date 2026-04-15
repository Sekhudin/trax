package generate

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var gRoutesCmd = docs.ApplyDocs(gRoutesDocs, &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
	},
})

func init() {
}
