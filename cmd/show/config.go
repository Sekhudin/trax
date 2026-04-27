package show

import (
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type showconfig struct {
	ctx app.Context
}

func NewConfigCmd(docs *doc.Docs, ctx app.Context) *cobra.Command {
	s := showconfig{ctx: ctx}
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return s.preRunE()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return s.runE(cmd)
		},
	})

	flags := cmd.Flags()
	flags.Bool("json", false, "output as json")

	return cmd
}

func (s *showconfig) preRunE() error {
	s.ctx.Out().Info("config", "show trax config\n")

	return nil
}

func (s *showconfig) runE(cmd *cobra.Command) error {
	asJSON, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}

	settings := viper.AllSettings()
	if asJSON {
		return s.ctx.Out().AsJSON(settings)
	}

	return s.ctx.Out().AsFlat("", settings)
}
