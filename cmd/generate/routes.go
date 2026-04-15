package generate

import (
	"fmt"

	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var RoutesCommand = docs.ApplyDocs(RoutesDocs, &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generating routes from")
	},
})

func init() {
}
