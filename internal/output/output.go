package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type Context struct {
	Writer io.Writer
}

type Notification struct {
	Level   string
	Scope   string
	Message string
}

func New(w io.Writer) *Context {
	return &Context{Writer: w}
}

func (c *Context) NotifyJSON(n Notification) {
	payload := map[string]string{
		"level":   n.Level,
		"scope":   n.Scope,
		"message": n.Message,
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(c.Writer, string(b))
}

func (c *Context) NotifyText(n Notification) {
	var format string
	switch n.Level {
	case "success":
		format = "✔ [%s] %s\n"
	case "warn":
		format = "  [%s] %s\n"
	default:
		format = "ℹ  [%s] %s\n"
	}

	fmt.Fprintf(c.Writer, format, n.Scope, n.Message)
}

func (c *Context) Success(scope, msg string) {
	notification := Notification{Level: "success", Scope: scope, Message: msg}
	if viper.GetBool("debug") {
		c.NotifyJSON(notification)
		return
	}

	c.NotifyText(notification)
}

func (c *Context) Info(scope, msg string) {
	notification := Notification{Level: "info", Scope: scope, Message: msg}
	if viper.GetBool("debug") {
		c.NotifyJSON(notification)
		return
	}

	c.NotifyText(notification)
}

func (c *Context) Warn(scope, msg string) {
	notification := Notification{Level: "warn", Scope: scope, Message: msg}
	if viper.GetBool("debug") {
		c.NotifyJSON(notification)
		return
	}

	c.NotifyText(notification)
}

func (c *Context) AsJSON(data map[string]any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return appErr.NewInternalError("config", "failed to marshal config to json", err)
	}

	if _, err := fmt.Fprintln(c.Writer, string(b)); err != nil {
		return appErr.NewIOError("stdout", "failed to write json output", err)
	}

	return nil
}

func (c *Context) AsFlat(prefix string, data map[string]any) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]

		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {

		case map[string]any:
			if err := c.AsFlat(key, val); err != nil {
				return err
			}

		case []any:
			for i, item := range val {
				itemKey := fmt.Sprintf("%s.%d", key, i)
				if err := c.printFlatValue(itemKey, item); err != nil {
					return err
				}
			}

		default:
			if err := c.printFlatValue(key, val); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Context) printFlatValue(key string, val any) error {
	switch v := val.(type) {

	case map[string]any:
		return c.AsFlat(key, v)

	default:
		if _, err := fmt.Fprintf(c.Writer, "%s = %v\n", key, v); err != nil {
			return appErr.NewIOError("stdout", "failed to write config output", err)
		}
	}

	return nil
}
