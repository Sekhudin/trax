package generate

import (
	"fmt"
	"path/filepath"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type generateconfig struct {
	ctx        app.Context
	cfgFormats map[string]struct{}
}

func NewConfigCmd(docs *doc.Docs, ctx app.Context) *cobra.Command {
	g := generateconfig{
		ctx: ctx,
		cfgFormats: map[string]struct{}{
			"json": {},
			"toml": {},
			"yaml": {},
			"yml":  {},
		},
	}

	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.runE(cmd)
		},
	})

	flags := cmd.Flags()
	flags.Bool("override", false, "overwrite existing config file")
	flags.StringP("format", "f", "toml", "config format")

	return cmd
}

func (g *generateconfig) runE(cmd *cobra.Command) error {
	flags := cmd.Flags()

	isOverride, err := flags.GetBool("override")
	if err != nil {
		return err
	}

	writeConfig := viper.SafeWriteConfigAs
	if isOverride {
		writeConfig = viper.WriteConfigAs
	}

	format, err := flags.GetString("format")
	if err != nil {
		return err
	}

	if !g.isValidConfigFormat(format) {
		return appErr.NewValidationError("config", fmt.Sprintf(
			"invalid value %q (allowed: %q, %q, or %q)", format, "json", "toml", "yaml",
		))
	}

	cfgFile := filepath.Clean(fmt.Sprintf("trax.%s", format))
	if err := writeConfig(cfgFile); err != nil {
		return appErr.NewConfigLoadError("config", fmt.Sprintf("failed generate config: %q", cfgFile), err)
	}

	g.ctx.Out().Success("routes", fmt.Sprintf("config written %s", g.ctx.Color().Green(cfgFile)))

	return nil
}

func (g *generateconfig) isValidConfigFormat(format string) bool {
	_, ok := g.cfgFormats[format]

	return ok
}
