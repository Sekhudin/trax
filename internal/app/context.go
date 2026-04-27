package app

import (
	"os"

	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/internal/runner"
	"github.com/spf13/cobra"
)

type Context interface {
	Color() output.Colorizer
	Out() output.Context
	Runner() runner.Runner

	ApplyOptions(cmd *cobra.Command, opt output.Options)
}

type context struct {
	color   output.Colorizer
	out     output.Context
	crunner runner.Runner
}

func New(opt output.Options) Context {
	out := output.New(os.Stdout, opt)

	return &context{
		out:     out,
		color:   out.Color(),
		crunner: runner.NewRunner(os.Stdout, os.Stderr),
	}
}

func (c *context) ApplyOptions(cmd *cobra.Command, opt output.Options) {
	out := output.New(os.Stdout, opt)

	c.out = out
	c.color = out.Color()
	c.crunner = runner.NewRunner(cmd.OutOrStdout(), cmd.ErrOrStderr())
}

func (c *context) Color() output.Colorizer {
	return c.color
}

func (c *context) Out() output.Context {
	return c.out
}

func (c *context) Runner() runner.Runner {
	return c.crunner
}
