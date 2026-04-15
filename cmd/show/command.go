package show

import (
	"trax/internal/docs"

	"github.com/spf13/cobra"
)

var Cmd = docs.ApplyDocs(sDocs, &cobra.Command{})

func init() {
	Cmd.AddCommand(sConfigCmd)
	Cmd.AddCommand(sRoutesCmd)
}
