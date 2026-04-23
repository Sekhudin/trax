package generate

import (
	"fmt"
	"path/filepath"

	"github.com/sekhudin/trax/internal/docs"
	"github.com/sekhudin/trax/internal/output"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type generateconfig struct {
	flags     *pflag.FlagSet
	out       *output.Context
	cfgFile   string
	cfgFormat map[string]struct{}
}

var (
	gc = generateconfig{
		cfgFormat: map[string]struct{}{
			"json": {},
			"toml": {},
			"yaml": {},
			"yml":  {},
		},
	}

	gcCommand = docs.ApplyDocs(&doc.config, &cobra.Command{
		RunE:     gc.runE,
		PostRunE: gc.postRunE,
	})
)

func init() {
	gc.flags = gcCommand.Flags()
	gc.out = output.New(gcCommand.OutOrStdout())

	gc.flags.Bool("override", false, "overwrite existing config file")
	gc.flags.StringP("format", "f", "toml", "config format")
}

func (g *generateconfig) runE(cmd *cobra.Command, args []string) error {
	isOverride, err := g.flags.GetBool("override")
	if err != nil {
		return err
	}

	writeConfig := viper.SafeWriteConfigAs
	if isOverride {
		writeConfig = viper.WriteConfigAs
	}

	format, err := g.flags.GetString("format")
	if err != nil {
		return err
	}

	if _, ok := gc.cfgFormat[format]; !ok {
		return appErr.NewValidationError("config", fmt.Sprintf(
			"invalid value %q (allowed: %q, %q, or %q)", format, "json", "toml", "yaml",
		))
	}

	g.cfgFile = filepath.Clean(fmt.Sprintf("trax.%s", format))
	if err := writeConfig(g.cfgFile); err != nil {
		return appErr.NewConfigLoadError("config", fmt.Sprintf("failed generate config: %q", g.cfgFile), err)
	}

	return nil
}

func (g *generateconfig) postRunE(cmd *cobra.Command, args []string) error {
	g.out.Success("routes", fmt.Sprintf("config written %s", output.Green(g.cfgFile)))

	return nil
}
