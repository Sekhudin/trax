package show

import (
	"fmt"

	"github.com/sekhudin/trax/internal/docs"
	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type showroutes struct {
	flags *pflag.FlagSet
	out   *output.Context
	cfg   *routes.Config
}

var (
	sr        = showroutes{}
	srCommand = docs.ApplyDocs(&doc.routes, &cobra.Command{
		PreRunE: sr.preRunE,
		RunE:    sr.runE,
	})
)

func init() {
	sr.flags = srCommand.Flags()
	sr.out = output.New(srCommand.OutOrStdout())

	sr.flags.StringP("strategy", "s", "", "route discovery strategy")
	sr.flags.StringP("root", "r", "", "project root directory used for route discovery")
	sr.flags.StringP("file", "f", "", "path to a route definition file")

	sr.flags.StringP("key", "k", "", "selector key")
	sr.flags.Bool("json", false, "output as json")

	srCommand.MarkFlagFilename("file", "yaml")
	srCommand.MarkFlagFilename("output", "ts", "js")
	srCommand.MarkFlagDirname("root")
}

func (s *showroutes) preRunE(cmd *cobra.Command, args []string) error {
	viper.BindPFlag("routes.strategy", s.flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", s.flags.Lookup("root"))
	viper.BindPFlag("routes.file", s.flags.Lookup("file"))

	cfg, err := routes.NewConfig()
	if err != nil {
		return err
	}

	s.cfg = cfg
	s.out.Info("routes", fmt.Sprintf("using %q strategy \n", s.cfg.Strategy))

	return nil
}

func (s *showroutes) runE(cmd *cobra.Command, args []string) error {
	key, err := s.flags.GetString("key")
	if err != nil {
		return err
	}

	asJSON, err := s.flags.GetBool("json")
	if err != nil {
		return err
	}

	selector, err := routes.Show(s.cfg)
	if err != nil {
		return err
	}

	val, err := selector(key)
	if err != nil {
		return err
	}

	if asJSON {
		s.out.AsJSON(val)
		return nil
	}

	s.out.AsFlat("", val)

	return nil
}
