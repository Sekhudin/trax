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

type generateroutes struct{}

var (
	gr        = generateroutes{}
	grCommand = docs.ApplyDocs(&doc.routes, &cobra.Command{
		PreRunE: gr.preRunE,
		RunE:    gr.runE,
	})
)

func init() {
	flags := grCommand.Flags()

	flags.StringP("strategy", "s", "", "route discovery strategy")
	flags.StringP("root", "r", "", "project root directory used for route discovery")
	flags.StringP("file", "f", "", "path to a route definition file")
	flags.StringP("output", "o", "", "output file path")

	grCommand.MarkFlagFilename("file", "yaml")
	grCommand.MarkFlagFilename("output", "ts", "js")
	grCommand.MarkFlagDirname("root")
}

func (*generateroutes) preRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))
	viper.BindPFlag("routes.output", flags.Lookup("output"))

	return nil
}

func (*generateroutes) runE(cmd *cobra.Command, args []string) error {
	out := output.New(cmd.OutOrStdout())

	cfg, err := routes.NewCfg()
	if err != nil {
		return appErr.NewValidationError("routes", err.Error())
	}

	key, _ := cmd.Flags().GetString("strategy")
	fmt.Println(key, cfg.Strategy)

	if cfg.Strategy == "file" {
		out.Info("routes", fmt.Sprintf("using route definition file: %q", cfg.File.Full))
		if err := routes.GenerateFromFile(cfg); err != nil {
			return appErr.NewIOError("routes", "generation failed", err)
		}

		return nil
	}

	out.Info("routes", fmt.Sprintf("discovering routes using strategy: %q", cfg.Strategy))
	if err := routes.GenerateFromDisc(cfg); err != nil {
		return appErr.NewIOError("routes", "generation failed", err)
	}

	return nil
}
