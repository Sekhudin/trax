package show

import (
	"trax/internal/docs"
	"trax/internal/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type showconfig struct{}

var (
	sc        = showconfig{}
	scCommand = docs.ApplyDocs(&doc.config, &cobra.Command{
		RunE: sc.runE,
	})
)

func init() {
	flags := scCommand.Flags()

	flags.Bool("json", false, "output as json")
}

func (*showconfig) runE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	out := output.New(cmd.OutOrStdout())

	asJSON, err := flags.GetBool("json")
	if err != nil {
		return appErr.NewFlagReadError("json", err)
	}

	settings := viper.AllSettings()
	if asJSON {
		return out.AsJSON(settings)
	}

	return out.AsFlat("", settings)
}
