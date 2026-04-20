package routes

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type treeselector func(selector string) (map[string]any, error)

type tree struct{}

func (*tree) build(rs []route) (map[string]*node, error) {
	tr := make(map[string]*node)

	for _, r := range rs {
		if err := r.insert(tr); err != nil {
			return nil, err
		}
	}

	return tr, nil
}

func (*tree) newSelector(tr map[string]any) (treeselector, error) {
	prefix := viper.GetString("routes.prefix")

	v := viper.New()

	v.SetDefault(prefix, tr[prefix])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", prefix, selector)
		} else {
			selector = prefix
		}

		selector = strings.TrimSuffix(selector, "$")
		selector = strings.TrimSuffix(selector, ".")
		selector = strings.TrimSuffix(selector, "?")
		selector = strings.ReplaceAll(selector, "?.", ".$")
		val := v.Get(selector)

		switch v := val.(type) {

		case nil:
			msg := fmt.Sprintf("%q not found", selector)

			return nil, appErr.NewValidationError("selector", msg)

		case map[string]any:
			return v, nil

		case string:
			return map[string]any{
				"path": v,
			}, nil

		default:
			msg := fmt.Sprintf("invalid type for selector %q: %T",
				selector, val,
			)

			return nil, appErr.NewValidationError("selector", msg)
		}
	}, nil
}

func (t *tree) toMap(nds map[string]*node) map[string]any {
	keys := make([]string, 0, len(nds))
	for k := range nds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make(map[string]any, len(nds))

	for _, k := range keys {
		n := nds[k]

		if n == nil {
			continue
		}

		m := make(map[string]any)

		if n.root != "" {
			m["root"] = n.root
		}

		if len(n.children) > 0 {
			maps.Copy(m, t.toMap(n.children))
		}

		result[k] = m
	}

	return result
}
