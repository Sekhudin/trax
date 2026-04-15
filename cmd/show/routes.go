package show

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var sRoutesCmd = docs.ApplyDocs(sRoutesDocs, &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
	},
})

func init() {
}
