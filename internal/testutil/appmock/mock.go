package appmock

import (
	"github.com/sekhudin/trax/internal/output"
	"github.com/sekhudin/trax/internal/runner"
	"github.com/sekhudin/trax/internal/testutil/outmock"
	"github.com/sekhudin/trax/internal/testutil/runnermock"
	"github.com/spf13/cobra"
)

type Context struct {
	ApplyOptionsCalled bool
	ColorCalled        bool
	OutCalled          bool
	RunnerCalled       bool

	OutMock    *outmock.Out
	ColorMock  *outmock.Color
	RunnerMock *runnermock.Runner
}

func NewContext() *Context {
	return &Context{
		OutMock:    &outmock.Out{},
		ColorMock:  &outmock.Color{},
		RunnerMock: &runnermock.Runner{},
	}
}

func (c *Context) Reset() {
	c.ApplyOptionsCalled = false
	c.ColorCalled = false
	c.OutCalled = false
	c.RunnerCalled = false

	c.OutMock = &outmock.Out{}
	c.ColorMock = &outmock.Color{}
	c.RunnerMock = &runnermock.Runner{}
}

func (c *Context) ApplyOptions(cmd *cobra.Command, opt output.Options) {
	c.ApplyOptionsCalled = true
}

func (c *Context) Color() output.Colorizer {
	c.ColorCalled = true
	return c.ColorMock
}

func (c *Context) Out() output.Context {
	c.OutCalled = true
	return c.OutMock
}

func (c *Context) Runner() runner.Runner {
	c.RunnerCalled = true
	return c.RunnerMock
}
