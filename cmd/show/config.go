package show

import (
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigCtx interface {
	PreRunE() error
	RunE(cmd *cobra.Command) error
}

type configctx struct {
	ctx app.Context
}

func NewConfigCmd(docs *doc.Docs, c ConfigCtx) *cobra.Command {
	cmd := doc.Apply(docs, &cobra.Command{
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.PreRunE()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return c.RunE(cmd)
		},
	})

	flags := cmd.Flags()
	flags.Bool("json", false, "output as json")

	return cmd
}

func NewConfigCtx(ctx app.Context) ConfigCtx {
	return &configctx{
		ctx: ctx,
	}
}

func (c *configctx) PreRunE() error {
	c.ctx.Out().Info("config", "show trax config\n")

	return nil
}

func (c *configctx) RunE(cmd *cobra.Command) error {
	asJSON, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}

	settings := viper.AllSettings()
	if asJSON {
		return c.ctx.Out().AsJSON(settings)
	}

	return c.ctx.Out().AsFlat("", settings)
}
