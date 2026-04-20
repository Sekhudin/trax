package show

import (
	"fmt"

	"trax/internal/docs"
	"trax/internal/output"
	"trax/modules/routes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type showroutes struct{}

var (
	sr        = showroutes{}
	srCommand = docs.ApplyDocs(&doc.routes, &cobra.Command{
		PreRunE: sr.preRunE,
		RunE:    sr.runE,
	})
)

func init() {
	flags := srCommand.Flags()

	flags.StringP("strategy", "s", "", "route discovery strategy")
	flags.StringP("root", "r", "", "project root directory used for route discovery")
	flags.StringP("file", "f", "", "path to a route definition file")

	flags.StringP("key", "k", "", "selector key")
	flags.Bool("json", false, "output as json")

	srCommand.MarkFlagFilename("file", "yaml")
	srCommand.MarkFlagFilename("output", "ts", "js")
	srCommand.MarkFlagDirname("root")
}

func (*showroutes) preRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	viper.BindPFlag("routes.strategy", flags.Lookup("strategy"))
	viper.BindPFlag("routes.root", flags.Lookup("root"))
	viper.BindPFlag("routes.file", flags.Lookup("file"))
	return nil
}

func (*showroutes) runE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	out := output.New(cmd.OutOrStdout())

	key, err := flags.GetString("key")
	if err != nil {
		return appErr.NewFlagReadError("key", err)
	}

	asJSON, err := flags.GetBool("json")
	if err != nil {
		return appErr.NewFlagReadError("json", err)
	}

	cfg, err := routes.NewCfg()
	if err != nil {
		return appErr.NewValidationError("routes", err.Error())
	}

	var tSelector routes.TreeSelector

	switch cfg.Strategy {
	case "file":
		rSelector, err := routes.ShowFromFile(cfg)
		if err != nil {
			return appErr.NewValidationError("routes", err.Error())
		}

		tSelector = rSelector
		out.Info("routes", fmt.Sprintf("using route definition file: %q\n", cfg.File.Full))

	default:
		rSelector, err := routes.ShowFromDisc(cfg)
		if err != nil {
			return appErr.NewValidationError("routes", err.Error())
		}

		tSelector = rSelector
		out.Info("routes", fmt.Sprintf("discovering routes using strategy: %q\n", cfg.Strategy))
	}

	val, err := tSelector(key)
	if err != nil {
		return appErr.NewIOError("routes", "path not found", err)
	}

	if asJSON {
		out.AsJSON(val)
		return nil
	}

	out.AsFlat("", val)
	return nil
}
