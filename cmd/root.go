package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"trax/cmd/generate"
	"trax/cmd/show"
	"trax/internal/bootstrap"
	"trax/internal/clierror"
	"trax/internal/docs"

	appErr "trax/internal/errors"
)

type trax struct {
	doc    docs.Docs
	flags  *pflag.FlagSet
	pFlags *pflag.FlagSet
}

var (
	tx = trax{
		doc: docs.Docs{
			Use:     "trax",
			Version: "0.0.1",
			Short:   "Powering TypeScript project workflows",
			Long: docs.Paragraph(
				"Trax is a CLI tool for automating TypeScript project workflows.",
			),
		},
	}

	command = docs.ApplyDocs(&tx.doc, &cobra.Command{
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: tx.persistentPreRunE,
	})
)

func init() {
	tx.flags = command.Flags()
	tx.pFlags = command.PersistentFlags()

	tx.pFlags.BoolP("debug", "d", false, "show debug info")
	tx.pFlags.Bool("no-color", false, "disable color")
	tx.pFlags.String("config", "", "path to config file")

	viper.BindPFlag("debug", tx.pFlags.Lookup("debug"))
	viper.BindPFlag("no-color", tx.pFlags.Lookup("no-color"))

	command.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return appErr.NewValidationError("flag", err.Error())
	})

	command.AddCommand(generate.Command, show.Command)
}

func (t *trax) persistentPreRunE(cmd *cobra.Command, args []string) error {
	cfgFile, err := t.flags.GetString("config")
	if err != nil {
		return appErr.NewFlagReadError("config", err)
	}
	return bootstrap.LoadConfig(cfgFile)
}

func Execute() {
	if cmd, err := command.ExecuteC(); err != nil {
		cErr := clierror.New(cmd.ErrOrStderr())

		cErr.Print(err)
		os.Exit(cErr.ExitCode(err))
	}
}
