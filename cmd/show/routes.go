package show

import (
	"fmt"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RoutesCtx interface {
	PreRunE(cmd *cobra.Command) error
	RunE(cmd *cobra.Command) error
}

type routesctx struct {
	ctx          app.Context
	cfg          *routes.Config
	config       func() config.Config
	routeConfig  func(*config.RoutesConfig) routes.RoutesConfig
	routeBuilder func(*routes.Config) routes.Builder
}

func NewRoutesCmd(docs *doc.Docs, c RoutesCtx) *cobra.Command {
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.PreRunE(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return c.RunE(cmd)
		},
	})

	flags := cmd.Flags()

	flags.StringP("strategy", "s", "", "route discovery strategy")
	flags.StringP("root", "r", "", "project root directory used for route discovery")
	flags.StringP("file", "f", "", "path to a route definition file")

	flags.StringP("key", "k", "", "selector key")
	flags.Bool("json", false, "output as json")

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
	}
}

func (c *routesctx) PreRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))

	cfg, err := c.routeConfig(c.config().Routes()).Load()
	if err != nil {
		fmt.Println("BELUM DI CHECK KE SINI 1")
		return err
	}

	c.cfg = cfg
	c.ctx.Out().Info("routes", fmt.Sprintf("using %q strategy \n", cfg.Strategy))

	return nil
}

func (c *routesctx) RunE(cmd *cobra.Command) error {
	flags := cmd.Flags()

	key, err := flags.GetString("key")
	if err != nil {
		fmt.Println("BELUM DI CHECK KE SINI 2")
		return err
	}

	asJSON, err := flags.GetBool("json")
	if err != nil {
		fmt.Println("BELUM DI CHECK KE SINI 3")
		return err
	}

	r, err := c.routeBuilder(c.cfg).Build()
	if err != nil {
		fmt.Println("BELUM DI CHECK KE SINI 4")
		return err
	}

	val, err := r.Select(key)
	if err != nil {
		fmt.Println("BELUM DI CHECK KE SINI 5")
		return err
	}

	if asJSON {
		c.ctx.Out().AsJSON(val)
		fmt.Println("BELUM DI CHECK KE SINI 6")
	} else {
		fmt.Println("BELUM DI CHECK KE SINI 7")
		c.ctx.Out().AsFlat("", val)
	}

	return nil
}
