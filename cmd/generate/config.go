package generate

import (
	"fmt"
	"path/filepath"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type ConfigCtx interface {
	RunE(cmd *cobra.Command) error
}

type configctx struct {
	ctx         app.Context
	cfgFormats  map[string]struct{}
	cfgFilePath func(string) string
	cfgWriter   func(string, bool) config.Writer
}

func NewConfigCmd(docs *doc.Docs, c ConfigCtx) *cobra.Command {
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.RunE(cmd)
		},
	})

	flags := cmd.Flags()
	flags.Bool("override", false, "overwrite existing config file")
	flags.StringP("format", "f", "toml", "config format")

	return cmd
}

func NewConfigCtx(ctx app.Context) ConfigCtx {
	return &configctx{
		ctx: ctx,
		cfgFormats: map[string]struct{}{
			"json": {},
			"toml": {},
			"yaml": {},
			"yml":  {},
		},
		cfgFilePath: func(s string) string {
			return filepath.Clean("trax." + s)
		},
		cfgWriter: func(s string, o bool) config.Writer {
			return config.NewWriter(s, o, viper.SafeWriteConfigAs, viper.WriteConfigAs)
		},
	}
}

func (c *configctx) RunE(cmd *cobra.Command) error {
	flags := cmd.Flags()

	override, err := flags.GetBool("override")
	if err != nil {
		return err
	}

	format, err := flags.GetString("format")
	if err != nil {
		return err
	}

	if _, ok := c.cfgFormats[format]; !ok {
		return appErr.NewValidationError("config", fmt.Sprintf(
			"invalid value %q (allowed: %q, %q, or %q)", format, "json", "toml", "yaml",
		))
	}

	w := c.cfgWriter(c.cfgFilePath(format), override)
	if err := w.Write(); err != nil {
		return appErr.NewConfigLoadError("config", "failed generate config", err)
	}

	c.ctx.Out().Success("routes", fmt.Sprintf("config written %s", c.ctx.Color().Green(w.File())))

	return nil
}
