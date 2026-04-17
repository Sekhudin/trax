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

var root = docs.ApplyDocs(rootDocs, &cobra.Command{
	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRunE: rootPersistentPreRunE,
})

func init() {
	pFlags := root.PersistentFlags()

	pFlags.BoolP("debug", "d", false, "show debug info")
	pFlags.Bool("no-color", false, "disable color")
	pFlags.String("config", "", "path to config file")

	viper.BindPFlag("debug", pFlags.Lookup("debug"))
	viper.BindPFlag("no-color", pFlags.Lookup("no-color"))

	root.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return appErr.NewValidationError("flag", err.Error())
	})

	root.AddCommand(generate.Cmd, show.Cmd)
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	cfgFile, err := flags.GetString("config")
	if err != nil {
		return appErr.NewFlagReadError("config", err)
	}
	return bootstrap.LoadConfig(cfgFile)
}

func Execute() {
	if cmd, err := root.ExecuteC(); err != nil {
		cErr := clierror.New(cmd.ErrOrStderr())

		cErr.Print(err)
		os.Exit(cErr.ExitCode(err))
	}
}
