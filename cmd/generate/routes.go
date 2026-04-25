package generate

import (
	"fmt"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type generateroutes struct {
	ctx *app.Context
	cfg *routes.Config
}

func NewRoutesCmd(docs *doc.Docs, ctx *app.Context) *cobra.Command {
	g := generateroutes{ctx: ctx}
	cmd := doc.Apply(docs, &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return g.preRunE(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return g.runE()
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			return g.postRunE(cmd)
		},
	})

	flags := cmd.Flags()
	flags.StringP("strategy", "s", "", "route discovery strategy")
	flags.StringP("root", "r", "", "project root directory used for route discovery")
	flags.StringP("file", "f", "", "path to a route definition file")
	flags.StringP("output", "o", "", "output file path")

	flags.StringP("formatter", "t", "", "specify code formatter to use")
	flags.BoolP("no-deps", "x", false, "do not include dependencies in the output")
	flags.BoolP("no-format", "n", false, "disable automatic code formatting")

	cmd.MarkFlagFilename("file", "yaml")
	cmd.MarkFlagFilename("output", "ts", "js")
	cmd.MarkFlagDirname("root")

	return cmd
}

func (g *generateroutes) preRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))
	viper.BindPFlag("routes.output", flags.Lookup("output"))
	viper.BindPFlag("routes.no-deps", flags.Lookup("no-deps"))

	viper.BindPFlag("formatter", flags.Lookup("formatter"))
	cfg, err := routes.NewConfig()
	if err != nil {
		return err
	}

	g.setCfg(cfg)

	return nil
}

func (g *generateroutes) runE() error {
	g.ctx.Out.Info("routes", fmt.Sprintf("using %q strategy (no-deps: %v)\n", g.cfg.Strategy, g.ctx.Color.Blue(g.cfg.NoDeps)))
	if err := routes.Generate(g.cfg); err != nil {
		return err
	}

	return nil
}

func (g *generateroutes) postRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	noformat, err := flags.GetBool("no-format")
	if err != nil {
		return err
	}

	if !noformat {
		f := viper.GetString("formatter")
		sf := viper.GetStringMap(fmt.Sprintf("formatters.%s", f))

		if err := g.ctx.Runner.Run(sf); err != nil {
			return err
		}
	}

	g.ctx.Out.Success("routes", fmt.Sprintf("routes written %s", g.ctx.Color.Green(g.cfg.Output.Full)))

	return nil
}

func (g *generateroutes) setCfg(cfg *routes.Config) {
	g.cfg = cfg
}
