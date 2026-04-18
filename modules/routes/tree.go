package routes

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type TreeSelector func(selector string) (map[string]any, error)

func buildTree(rs []route) (map[string]*node, error) {
	tree := make(map[string]*node)

	for _, r := range rs {
		if err := r.insert(tree); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func newTreeSelector(tree map[string]any) (TreeSelector, error) {
	prefix := viper.GetString("routes.prefix")

	v := viper.New()

	v.SetDefault(prefix, tree[prefix])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", prefix, selector)
		} else {
			selector = prefix
		}

		selector = strings.TrimSuffix(selector, ".")
		selector = strings.ReplaceAll(selector, "?", "$")
		val := v.Get(selector)

		switch v := val.(type) {

		case nil:
			return nil, fmt.Errorf("%q not found", selector)

		case map[string]any:
			return v, nil

		case string:
			return map[string]any{
				"path": v,
			}, nil

		default:
			return nil, fmt.Errorf(
				"invalid type for selector %q: %T",
				selector,
				val,
			)
		}
	}, nil
}

func toMap(nds map[string]*node) map[string]any {
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
			maps.Copy(m, toMap(n.children))
		}

		result[k] = m
	}

	return result
}
