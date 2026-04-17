package clierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
	"trax/internal/output"
)

type Context struct {
	Writer io.Writer
}

func New(w io.Writer) *Context {
	return &Context{Writer: w}
}

func (c *Context) PrintText(err error) {
	switch e := err.(type) {
	case *appErr.CoreError:
		fmt.Fprintf(c.Writer, "%s %s:", output.IconError, output.Bold(e.Code))

		if e.Scope != "" {
			fmt.Fprintf(c.Writer, " (%s)", output.Red(e.Scope))
		}

		fmt.Fprintf(c.Writer, " %s\n", e.Message)

		if e.Err != nil {
			fmt.Fprintf(c.Writer, "    %s %v\n", output.IconDetail, e.Err)
		}

	default:
		fmt.Fprintf(c.Writer, "%s %v\n", output.IconError, err)
	}
}

func (c *Context) PrintJSON(err error) {
	payload := map[string]any{
		"error": err.Error(),
	}

	if e, ok := err.(*appErr.CoreError); ok {
		payload["code"] = e.Code
		payload["scope"] = e.Scope
		payload["message"] = e.Message

		if e.Err != nil {
			payload["cause"] = e.Err.Error()
		}
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(c.Writer, string(b))
}

func (c *Context) Print(err error) {
	if err == nil {
		return
	}

	if viper.GetBool("debug") {
		c.PrintJSON(err)
		return
	}

	c.PrintText(err)
}

func (c *Context) ExitCode(err error) int {
	var ce *appErr.CoreError
	if errors.As(err, &ce) {
		switch ce.Code {
		case appErr.ErrValidation:
			return 2
		case appErr.ErrConfigNotFound:
			return 3
		case appErr.ErrConfigLoad:
			return 4
		default:
			return 1
		}
	}
	return 1
}
