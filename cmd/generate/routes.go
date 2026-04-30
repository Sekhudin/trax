package generate

import (
	"fmt"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/fs"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RoutesCtx interface {
	PreRunE(cmd *cobra.Command) error
	RunE() error
	PostRunE(cmd *cobra.Command) error
}

type routesctx struct {
	ctx            app.Context
	cfg            *routes.Config
	config         func() config.Config
	routeConfig    func(*config.RoutesConfig) routes.RoutesConfig
	routeBuilder   func(*routes.Config) routes.Builder
	routeTemplate  func([]routes.Route, routes.TreeSelector, *routes.Config) routes.Template
	routeGenerator func(routes.Template) routes.Generator
}

func NewRoutesCmd(docs *doc.Docs, c RoutesCtx) *cobra.Command {
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.PreRunE(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return c.RunE()
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			return c.PostRunE(cmd)
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

func NewRoutesCtx(ctx app.Context) RoutesCtx {
	return &routesctx{
		ctx: ctx,
		cfg: nil,
		config: func() config.Config {
			return config.New()
		},
		routeConfig: func(c *config.RoutesConfig) routes.RoutesConfig {
			return routes.NewConfig(c)
		},
		routeBuilder: func(c *routes.Config) routes.Builder {
			return routes.NewBuilder(c)
		},
		routeTemplate: func(r []routes.Route, ts routes.TreeSelector, c *routes.Config) routes.Template {
			return routes.NewTemplate(routes.TemplateDeps{
				Routes:   r,
				Selector: ts,
				Cfg:      c,
			})
		},
		routeGenerator: func(t routes.Template) routes.Generator {
			return routes.NewGenerator(fs.NewOSWriter(), t)
		},
	}
}

func (c *routesctx) PreRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))
	viper.BindPFlag("routes.output", flags.Lookup("output"))
	viper.BindPFlag("routes.no-deps", flags.Lookup("no-deps"))

	viper.BindPFlag("formatter", flags.Lookup("formatter"))

	cfg, err := c.routeConfig(c.config().Routes()).Load()
	if err != nil {
		return err
	}

	c.cfg = cfg
	c.ctx.Out().Info("routes", fmt.Sprintf("using %q strategy (no-deps: %v)\n", cfg.Strategy, c.ctx.Color().Blue(cfg.NoDeps)))

	return nil
}

func (c *routesctx) RunE() error {
	r, err := c.routeBuilder(c.cfg).Build()
	if err != nil {
		return err
	}

	return c.routeGenerator(c.routeTemplate(
		r.Routes(),
		r.Selector(),
		c.cfg,
	)).Generate(c.cfg.Output.Full)
}

func (c *routesctx) PostRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	noformat, err := flags.GetBool("no-format")
	if err != nil {
		return err
	}

	if !noformat {
		f := viper.GetString("formatter")
		sf := viper.GetStringMap(fmt.Sprintf("formatters.%s", f))

		if err := c.ctx.Runner().Run(sf); err != nil {
			return err
		}
	}

	c.ctx.Out().Success("routes", fmt.Sprintf("routes written %s", c.ctx.Color().Green(c.cfg.Output.Filename)))

	return nil
}
