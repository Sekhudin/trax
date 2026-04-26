package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sekhudin/trax/cmd/generate"
	"github.com/sekhudin/trax/cmd/show"
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/bootstrap"
	"github.com/sekhudin/trax/internal/clierror"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/output"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type trax struct {
	ctx *app.Context
}

var Version = ""

func New() *cobra.Command {
	t := trax{ctx: app.New(output.Options{})}

	cmd := &cobra.Command{
		Use:     "trax",
		Version: app.Version(Version),
		Short:   "Powering TypeScript project workflows",
		Long: doc.Paragraph(
			"Trax is a CLI tool for automating TypeScript project workflows.",
		),
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return t.persistentPreRunE(cmd)
		},
	}

	pFlags := cmd.PersistentFlags()
	pFlags.BoolP("debug", "d", false, "show debug info")
	pFlags.Bool("no-color", false, "disable color")
	pFlags.String("config", "", "path to config file")

	viper.BindPFlag("debug", pFlags.Lookup("debug"))
	viper.BindPFlag("no-color", pFlags.Lookup("no-color"))

	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return appErr.NewValidationError("flag", err.Error())
	})

	cmd.AddCommand(generate.New(t.ctx), show.New(t.ctx))

	return cmd
}

func (t *trax) persistentPreRunE(cmd *cobra.Command) error {
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return appErr.NewFlagReadError("config", err)
	}

	t.ctx.ApplyOptions(cmd, output.Options{
		Debug:   viper.GetBool("debug"),
		NoColor: viper.GetBool("no-color"),
	})

	return bootstrap.LoadConfig(cfgFile)
}

func Execute() {
	command := New()

	if cmd, err := command.ExecuteC(); err != nil {
		cErr := clierror.New(output.New(cmd.OutOrStdout(), output.Options{
			Debug:   viper.GetBool("debug"),
			NoColor: viper.GetBool("no-color"),
		}))

		cErr.Print(err)
		os.Exit(cErr.ExitCode(err))
	}
}
