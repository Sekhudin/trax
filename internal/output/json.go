package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	appErr "github.com/sekhudin/trax/internal/errors"
)

func (c *context) AsJSON(data map[string]any) error {
	normalized := c.normalizeMap(data)

	b, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return appErr.NewInternalError("json", "failed to marshal data to json", err)
	}

	if _, err := fmt.Fprintln(c.w, string(b)); err != nil {
		return appErr.NewIOError("stdout", "failed to write json output", err)
	}

	return nil
}

func (c *context) normalizeMap(m map[string]any) map[string]any {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")

	for i, k := range keys {
		v := c.normalizeValue(m[k])

		keyJSON, _ := json.Marshal(k)
		valJSON, _ := json.Marshal(v)

		buf.Write(keyJSON)
		buf.WriteString(":")
		buf.Write(valJSON)

		if i < len(keys)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")

	var out map[string]any
	_ = json.Unmarshal(buf.Bytes(), &out)
	return out
}

func (c *context) normalizeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return c.normalizeMap(val)

	case []any:
		arr := make([]any, len(val))
		for i, item := range val {
			arr[i] = c.normalizeValue(item)
		}
		return arr

	default:
		return val
	}
}
