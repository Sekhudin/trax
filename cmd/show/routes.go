package show

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
		return appErr.NewValidationError("routes", err.Error())
	}

	s.cfg = cfg
	s.out.Info("routes", fmt.Sprintf("showing using %q strategy \n", s.cfg.Strategy))

	return nil
}

func (s *showroutes) runE(cmd *cobra.Command, args []string) error {
	key, err := s.flags.GetString("key")
	if err != nil {
		return appErr.NewFlagReadError("key", err)
	}

	asJSON, err := s.flags.GetBool("json")
	if err != nil {
		return appErr.NewFlagReadError("json", err)
	}

	selector, err := routes.Show(s.cfg)
	if err != nil {
		return appErr.NewValidationError("routes", err.Error())
	}

	val, err := selector(key)
	if err != nil {
		return appErr.NewIOError("routes", "path not found", err)
	}

	if asJSON {
		s.out.AsJSON(val)
		return nil
	}

	s.out.AsFlat("", val)

	return nil
}
