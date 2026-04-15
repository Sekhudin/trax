package show

import (
	"fmt"

	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var RoutesCommand = docs.ApplyDocs(RoutesDocs, &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Showing routes from")
	},
})

func init() {
}
