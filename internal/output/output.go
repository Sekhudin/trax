package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/viper"
)

type Context struct {
	Writer io.Writer
}

func New(w io.Writer) *Context {
	return &Context{Writer: w}
}

func (c *Context) Success(scope, msg string) {
	if viper.GetBool("debug") {
		printJson(c.Writer, "success", scope, msg)
		return
	}

	printText(c.Writer, "✔ [%s] %s\n", scope, msg)
}

func (c *Context) Info(scope, msg string) {
	if viper.GetBool("debug") {
		printJson(c.Writer, "info", scope, msg)
		return
	}

	printText(c.Writer, "ℹ  [%s] %s\n", scope, msg)
}

func (c *Context) Warn(scope, msg string) {
	if viper.GetBool("debug") {
		printJson(c.Writer, "warn", scope, msg)
		return
	}

	printText(c.Writer, "  [%s] %s\n", scope, msg)
}

func printText(w io.Writer, format, scope, msg string) {
	fmt.Fprintf(w, format, scope, msg)
}

func printJson(w io.Writer, level, scope, msg string) {
	payload := map[string]string{
		"level":   level,
		"scope":   scope,
		"message": msg,
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(w, "\n"+string(b))
}
