package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"trax/cmd/generate"
	"trax/cmd/show"
	"trax/internal/bootstrap"
	"trax/internal/clierror"
	"trax/internal/docs"

	appErr "trax/internal/errors"
)

type trax struct {
	doc docs.Docs
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
	pFlags := command.PersistentFlags()

	pFlags.BoolP("debug", "d", false, "show debug info")
	pFlags.Bool("no-color", false, "disable color")
	pFlags.String("config", "", "path to config file")

	command.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return appErr.NewValidationError("flag", err.Error())
	})

	command.AddCommand(generate.Command, show.Command)
}

func (*trax) persistentPreRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	pFlags := cmd.PersistentFlags()

	viper.BindPFlag("debug", pFlags.Lookup("debug"))
	viper.BindPFlag("no-color", pFlags.Lookup("no-color"))

	cfgFile, err := flags.GetString("config")
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
