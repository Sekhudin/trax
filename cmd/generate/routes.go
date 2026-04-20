package generate

import (
	"fmt"

	"trax/internal/docs"
	"trax/internal/output"
	"trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type generateroutes struct {
	flags *pflag.FlagSet
	out   *output.Context
	cfg   *routes.Config
}

var (
	gr        = generateroutes{}
	grCommand = docs.ApplyDocs(&doc.routes, &cobra.Command{
		PreRunE:  gr.preRunE,
		RunE:     gr.runE,
		PostRunE: gr.postRunE,
	})
)

func init() {
	gr.flags = grCommand.Flags()
	gr.out = output.New(grCommand.OutOrStdout())

	gr.flags.StringP("strategy", "s", "", "route discovery strategy")
	gr.flags.StringP("root", "r", "", "project root directory used for route discovery")
	gr.flags.StringP("file", "f", "", "path to a route definition file")
	gr.flags.StringP("output", "o", "", "output file path")

	grCommand.MarkFlagFilename("file", "yaml")
	grCommand.MarkFlagFilename("output", "ts", "js")
	grCommand.MarkFlagDirname("root")
}

func (g *generateroutes) preRunE(cmd *cobra.Command, args []string) error {
	viper.BindPFlag("routes.strategy", g.flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", g.flags.Lookup("root"))
	viper.BindPFlag("routes.file", g.flags.Lookup("file"))
	viper.BindPFlag("routes.output", g.flags.Lookup("output"))

	cfg, err := routes.NewConfig()
	if err != nil {
		return appErr.NewValidationError("routes", err.Error())
	}

	g.cfg = cfg
	g.out.Info("routes", fmt.Sprintf("generating using %q strategy\n", g.cfg.Strategy))

	return nil
}

func (g *generateroutes) runE(cmd *cobra.Command, args []string) error {
	if err := routes.Generate(g.cfg); err != nil {
		return err
	}

	return nil
}

func (g *generateroutes) postRunE(cmd *cobra.Command, args []string) error {
	g.out.Success("config", "routes written")

	return nil
}
