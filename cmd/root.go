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

	pFlags.String("config", "", "Path to config file")
	pFlags.BoolP("debug", "d", false, "Show debug info")

	viper.BindPFlag("debug", pFlags.Lookup("debug"))

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
		clierror.Print(cmd.ErrOrStderr(), err)
		os.Exit(clierror.ExitCode(err))
	}
}
