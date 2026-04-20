package show

import (
	"trax/internal/docs"
	"trax/internal/output"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type showconfig struct {
	flags *pflag.FlagSet
	out   *output.Context
}

var (
	sc        = showconfig{}
	scCommand = docs.ApplyDocs(&doc.config, &cobra.Command{
		PreRunE: sc.preRunE,
		RunE:    sc.runE,
	})
)

func init() {
	sc.flags = scCommand.Flags()
	sc.out = output.New(scCommand.OutOrStdout())

	sc.flags.Bool("json", false, "output as json")
}

func (s *showconfig) preRunE(cmd *cobra.Command, args []string) error {
	s.out.Info("config", "showing trax config\n")

	return nil
}

func (s *showconfig) runE(cmd *cobra.Command, args []string) error {
	asJSON, err := s.flags.GetBool("json")
	if err != nil {
		return appErr.NewFlagReadError("json", err)
	}

	settings := viper.AllSettings()
	if asJSON {
		return s.out.AsJSON(settings)
	}

	return s.out.AsFlat("", settings)
}
