package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sekhudin/trax/cmd/generate"
	"github.com/sekhudin/trax/cmd/show"
	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/bootstrap"
	"github.com/sekhudin/trax/internal/clierror"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/output"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type Docs struct {
	Root doc.Docs
}

type Dependencies struct {
	Docs           *Docs
	Ctx            Context
	NewGenerateCmd func(app.Context) *cobra.Command
	NewShowCmd     func(app.Context) *cobra.Command
}

type Context interface {
	FlagErrorFn(c *cobra.Command, err error) error
	PersistentPreRunE(cmd *cobra.Command) error
}

type context struct {
	ctx app.Context
}

var (
	Version = ""
	Command = func() *cobra.Command {
		return New(app.New(output.Options{}))
	}
	ErrorHanler = func(err error, h clierror.Handler) {
		h.Print(err)
		os.Exit(h.ExitCode(err))
	}
)

func New(ctx app.Context) *cobra.Command {
	return NewWithDependencies(ctx, DefaultDependencies(ctx))
}

func NewWithDependencies(ctx app.Context, d *Dependencies) *cobra.Command {
	cmd := doc.Apply(&d.Docs.Root, &cobra.Command{
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return d.Ctx.PersistentPreRunE(cmd)
		},
	})

	pFlags := cmd.PersistentFlags()
	pFlags.BoolP("debug", "d", false, "show debug info")
	pFlags.Bool("no-color", false, "disable color")
	pFlags.String("config", "", "path to config file")

	viper.BindPFlag("debug", pFlags.Lookup("debug"))
	viper.BindPFlag("no-color", pFlags.Lookup("no-color"))

	generateCmd := d.NewGenerateCmd(ctx)
	showCmd := d.NewShowCmd(ctx)

	cmd.SetFlagErrorFunc(d.Ctx.FlagErrorFn)
	cmd.AddCommand(generateCmd, showCmd)

	return cmd
}

func DefaultDependencies(ctx app.Context) *Dependencies {
	return &Dependencies{
		NewGenerateCmd: generate.New,
		NewShowCmd:     show.New,
		Ctx:            NewContext(ctx),
		Docs: &Docs{
			Root: doc.Docs{
				Use:     "trax",
				Version: app.Version(Version),
				Short:   "Powering TypeScript project workflows",
				Long: doc.Paragraph(
					"Trax is a CLI tool for automating TypeScript project workflows.",
				),
			},
		},
	}
}

func NewContext(ctx app.Context) Context {
	return &context{
		ctx: ctx,
	}
}

func (*context) FlagErrorFn(c *cobra.Command, err error) error {
	return appErr.NewValidationError("flag", err.Error())
}

func (c *context) PersistentPreRunE(cmd *cobra.Command) error {
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return appErr.NewFlagReadError("config", err)
	}

	c.ctx.ApplyOptions(cmd, output.Options{
		Debug:   viper.GetBool("debug"),
		NoColor: viper.GetBool("no-color"),
	})

	return bootstrap.LoadConfig(cfgFile)
}

func Execute() {
	command := Command()

	if cmd, err := command.ExecuteC(); err != nil {
		cErr := clierror.New(output.New(cmd.OutOrStdout(), output.Options{
			Debug:   viper.GetBool("debug"),
			NoColor: viper.GetBool("no-color"),
		}))

		ErrorHanler(err, cErr)
	}
}
