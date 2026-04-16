package generate

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var Cmd = docs.ApplyDocs(gDocs, &cobra.Command{})

func init() {
	Cmd.AddCommand(gConfigCmd, gRoutesCmd)
}
