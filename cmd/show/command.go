package show

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var Command = docs.ApplyDocs(Docs, &cobra.Command{})

func init() {
	Command.AddCommand(RoutesCommand)
}
