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
)

var command = docs.ApplyDocs(Docs, &cobra.Command{
	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRunE: persistentPreRunE,
})

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	return bootstrap.LoadConfig()
}

func Execute() {
	if err := command.Execute(); err != nil {
		clierror.Print(err)
		os.Exit(1)
	}
}

func init() {
	pFlags := command.PersistentFlags()

	pFlags.String("config", "", "Path to .trax.config file")
	pFlags.BoolP("debug", "d", false, "show debug info")

	viper.BindPFlag("config", pFlags.Lookup("config"))
	viper.BindPFlag("debug", pFlags.Lookup("debug"))

	command.AddCommand(generate.Command, show.Command)
}
