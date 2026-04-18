package generate

import (
	"fmt"

	"trax/internal/docs"
	"trax/internal/output"
	"trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

var gRoutesCmd = docs.ApplyDocs(gRoutesDocs, &cobra.Command{
	RunE: gRoutesRunE,
})

func init() {
	flags := gRoutesCmd.Flags()

	flags.StringP("strategy", "s", "", "route discovery strategy")
	flags.StringP("root", "r", "", "project root directory used for route discovery")
	flags.StringP("file", "f", "", "path to a route definition file")
	flags.StringP("output", "o", "", "output file path")

	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))
	viper.BindPFlag("routes.output", flags.Lookup("output"))

	gRoutesCmd.MarkFlagFilename("file", "yaml")
	gRoutesCmd.MarkFlagFilename("output", "ts", "js")
	gRoutesCmd.MarkFlagDirname("root")
}

func gRoutesRunE(cmd *cobra.Command, args []string) error {
	out := output.New(cmd.OutOrStdout())

	cfg, err := routes.LoadConfig()
	if err != nil {
		return appErr.NewValidationError("routes", err.Error())
	}

	if cfg.Strategy == "file" {
		if err := routes.GenerateFromFile(cfg); err != nil {
			return appErr.NewIOError("routes", "", err)
		}

		out.Info("routes", fmt.Sprintf("using route definition file: %q", cfg.File.Full))
		return nil
	}

	if err := routes.GenerateFromDisc(cfg); err != nil {
		return appErr.NewIOError("routes", "", err)
	}

	out.Info("routes", fmt.Sprintf("discovering routes using strategy: %q", cfg.Strategy))
	return nil
}
