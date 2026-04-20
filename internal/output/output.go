package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type context struct {
	Writer io.Writer
}

type notification struct {
	Level   string
	Scope   string
	Message string
}

func New(w io.Writer) *context {
	return &context{Writer: w}
}

func (c *context) notifyJSON(n notification) {
	payload := map[string]string{
		"level":   n.Level,
		"scope":   n.Scope,
		"message": n.Message,
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(c.Writer, string(b))
}

func (c *context) notifyText(n notification) {
	var icon string
	switch n.Level {

	case "success":
		icon = IconSuccess

	case "warn":
		icon = IconWarn

	default:
		icon = IconInfo
	}

	fmt.Fprintf(c.Writer, "%s (%s) %s\n", icon, n.Scope, n.Message)
}

func (c *context) Success(scope, msg string) {
	notification := notification{Level: "success", Scope: Green(scope), Message: msg}
	if viper.GetBool("debug") {
		c.notifyJSON(notification)
		return
	}

	c.notifyText(notification)
}

func (c *context) Info(scope, msg string) {
	notification := notification{Level: "info", Scope: Blue(scope), Message: msg}
	if viper.GetBool("debug") {
		c.notifyJSON(notification)
		return
	}

	c.notifyText(notification)
}

func (c *context) Warn(scope, msg string) {
	notification := notification{Level: "warn", Scope: Yellow(scope), Message: msg}
	if viper.GetBool("debug") {
		c.notifyJSON(notification)
		return
	}

	c.notifyText(notification)
}

func (c *context) AsJSON(data map[string]any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return appErr.NewInternalError("config", "failed to marshal data to json", err)
	}

	if _, err := fmt.Fprintln(c.Writer, string(b)); err != nil {
		return appErr.NewIOError("stdout", "failed to show output", err)
	}

	return nil
}

func (c *context) AsFlat(prefix string, data map[string]any) error {
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

func (c *context) printFlatValue(key string, val any) error {
	switch v := val.(type) {

	case map[string]any:
		return c.AsFlat(key, v)

	default:
		if _, err := fmt.Fprintf(c.Writer, "%s = %v\n", key, v); err != nil {
			return err
		}
	}

	return nil
}
