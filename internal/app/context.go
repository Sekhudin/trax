package app

import (
	"os"
	"runtime/debug"

	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/internal/runner"
)

type Context struct {
	Version string
	Color   output.Colorizer
	Out     *output.Context
	OutOpt  output.Options
	Runner  runner.Runner
}

func New(v string) *Context {
	outOpt := output.Options{}
	return &Context{
		Version: version(v),
		Color:   output.NewColorizer(outOpt.NoColor),
		Out:     output.New(os.Stdout, outOpt),
		OutOpt:  outOpt,
		Runner:  runner.NewRunner(os.Stdout, os.Stderr),
	}
}

func version(v string) string {
	if v != "" {
		return v
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}

func (c *Context) SetOutOptDebug(val bool) {
	c.OutOpt.Debug = val
}

func (c *Context) SetOutOptNoColor(val bool) {
	c.OutOpt.NoColor = val
	c.Color.NoColor = val
}
