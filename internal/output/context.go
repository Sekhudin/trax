package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Level int

const (
	LevelInfo Level = iota
	LevelSuccess
	LevelWarn
	LevelError
)

type Options struct {
	Debug   bool
	JSON    bool
	Quiet   bool
	NoColor bool
}

type Context struct {
	w     io.Writer
	opt   Options
	color Colorizer
}

func New(w io.Writer, opt Options) *Context {
	return &Context{
		w:     w,
		opt:   opt,
		color: NewColorizer(opt.NoColor),
	}
}

type notification struct {
	Level   string `json:"level"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
}

func (c *Context) notify(level Level, scope, msg string) {
	if c.opt.Quiet {
		return
	}

	if c.opt.JSON || c.opt.Debug {
		c.notifyJSON(level, scope, msg)
		return
	}

	c.notifyText(level, scope, msg)
}

func (c *Context) notifyJSON(level Level, scope, msg string) {
	n := notification{
		Level:   level.String(),
		Scope:   scope,
		Message: msg,
	}

	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		fmt.Fprintf(c.w, `{"level":"error","message":"failed to marshal json"}\n`)
		return
	}

	fmt.Fprintln(c.w, string(b))
}

func (c *Context) notifyText(level Level, scope, msg string) {
	icon := c.icon(level)
	scope = c.colorScope(level, scope)

	fmt.Fprintf(c.w, "%s (%s) %s\n", icon, scope, msg)
}

func (c *Context) icon(level Level) string {
	switch level {
	case LevelSuccess:
		return c.color.Green("✔")
	case LevelWarn:
		return c.color.Yellow("⚠")
	case LevelError:
		return c.color.Red("✖")
	default:
		return c.color.Blue("ℹ")
	}
}

func (c *Context) colorScope(level Level, s string) string {
	switch level {
	case LevelSuccess:
		return c.color.Green(s)
	case LevelWarn:
		return c.color.Yellow(s)
	case LevelError:
		return c.color.Red(s)
	default:
		return c.color.Blue(s)
	}
}

func (l Level) String() string {
	switch l {
	case LevelSuccess:
		return "success"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

func (c *Context) Success(scope, msg string) {
	c.notify(LevelSuccess, scope, msg)
}

func (c *Context) Info(scope, msg string) {
	c.notify(LevelInfo, scope, msg)
}

func (c *Context) Warn(scope, msg string) {
	c.notify(LevelWarn, scope, msg)
}

func (c *Context) Error(scope, msg string) {
	c.notify(LevelError, scope, msg)
}
