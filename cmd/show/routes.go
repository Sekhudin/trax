package show

import (
	"fmt"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type showroutes struct {
	ctx app.Context
	cfg *routes.Config
}

func NewRoutesCmd(docs *doc.Docs, ctx app.Context) *cobra.Command {
	s := showroutes{ctx: ctx}
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return s.preRunE(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return s.runE(cmd)
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

func (s *showroutes) preRunE(cmd *cobra.Command) error {
	flags := cmd.Flags()
	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))

	cfg, err := routes.NewConfig()
	if err != nil {
		return err
	}

	s.setCfg(cfg)
	s.ctx.Out().Info("routes", fmt.Sprintf("using %q strategy \n", cfg.Strategy))

	return nil
}

func (s *showroutes) runE(cmd *cobra.Command) error {
	flags := cmd.Flags()

	key, err := flags.GetString("key")
	if err != nil {
		return err
	}

	asJSON, err := flags.GetBool("json")
	if err != nil {
		return err
	}

	b := routes.NewBuilder(s.cfg)
	r, err := b.Build()
	if err != nil {
		return err
	}

	val, err := r.Selector(key)
	if err != nil {
		return err
	}

	if asJSON {
		s.ctx.Out().AsJSON(val)
	} else {
		s.ctx.Out().AsFlat("", val)
	}

	return nil
}

func (s *showroutes) setCfg(cfg *routes.Config) {
	s.cfg = cfg
}
