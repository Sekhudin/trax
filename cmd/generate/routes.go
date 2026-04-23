package generate

import (
	"fmt"

	"github.com/sekhudin/trax/internal/docs"
	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/internal/runner"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type generateroutes struct {
	flags  *pflag.FlagSet
	out    *output.Context
	cfg    *routes.Config
	runner runner.Runner
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
	gr.runner = runner.NewRunner(grCommand.ErrOrStderr(), grCommand.ErrOrStderr())

	gr.flags.StringP("strategy", "s", "", "route discovery strategy")
	gr.flags.StringP("root", "r", "", "project root directory used for route discovery")
	gr.flags.StringP("file", "f", "", "path to a route definition file")
	gr.flags.StringP("output", "o", "", "output file path")

	gr.flags.StringP("formatter", "t", "", "specify code formatter to use")
	gr.flags.BoolP("no-deps", "x", false, "do not include dependencies in the output")
	gr.flags.BoolP("no-format", "n", false, "disable automatic code formatting")

	grCommand.MarkFlagFilename("file", "yaml")
	grCommand.MarkFlagFilename("output", "ts", "js")
	grCommand.MarkFlagDirname("root")
}

func (g *generateroutes) preRunE(cmd *cobra.Command, args []string) error {
	viper.BindPFlag("routes.strategy", g.flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", g.flags.Lookup("root"))
	viper.BindPFlag("routes.file", g.flags.Lookup("file"))
	viper.BindPFlag("routes.output", g.flags.Lookup("output"))
	viper.BindPFlag("routes.no-deps", g.flags.Lookup("no-deps"))

	viper.BindPFlag("formatter", g.flags.Lookup("formatter"))

	cfg, err := routes.NewConfig()
	if err != nil {
		return err
	}

	g.cfg = cfg
	g.out.Info("routes", fmt.Sprintf("using %q strategy (no-deps: %v)\n", g.cfg.Strategy, output.Blue(g.cfg.NoDeps)))

	return nil
}

func (g *generateroutes) runE(cmd *cobra.Command, args []string) error {
	if err := routes.Generate(g.cfg); err != nil {
		return err
	}

	return nil
}

func (g *generateroutes) postRunE(cmd *cobra.Command, args []string) error {
	noformat, err := g.flags.GetBool("no-format")
	if err != nil {
		return err
	}

	if !noformat {
		f := viper.GetString("formatter")
		sf := viper.GetStringMap(fmt.Sprintf("formatters.%s", f))

		if err := g.runner.Run(sf); err != nil {
			return err
		}
	}

	g.out.Success("routes", fmt.Sprintf("routes written %s", output.Green(g.cfg.Output.Full)))

	return nil
}
