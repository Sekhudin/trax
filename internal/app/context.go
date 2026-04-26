package app

import (
	"os"
	"runtime/debug"

	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/internal/runner"
	"github.com/spf13/cobra"
)

type Context struct {
	Color  *output.Colorizer
	Out    *output.Context
	Runner runner.Runner
}

func New(opt output.Options) *Context {
	return &Context{
		Color:  output.NewColorizer(opt.NoColor),
		Out:    output.New(os.Stdout, opt),
		Runner: runner.NewRunner(os.Stdout, os.Stderr),
	}
}

var readBuildInfo = debug.ReadBuildInfo

func Version(version string) string {
	if version != "" {
		return version
	}

	if info, ok := readBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}

func (c *Context) ApplyOptions(cmd *cobra.Command, opt output.Options) {
	c.Color = output.NewColorizer(opt.NoColor)
	c.Out = output.New(cmd.OutOrStdout(), opt)
	c.Runner = runner.NewRunner(cmd.OutOrStdout(), cmd.ErrOrStderr())
}
