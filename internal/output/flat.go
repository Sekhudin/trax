package output

import (
	"fmt"
	"sort"
)

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
	if _, err := fmt.Fprintf(c.w, "%s = %v\n", key, val); err != nil {
		return err
	}
	return nil
}
