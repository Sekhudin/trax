package routes

import (
	"fmt"
	"maps"
	"sort"

	"github.com/spf13/viper"
)

type treeSelector func(selector string) (map[string]any, error)

func buildTree(rs []route) (map[string]*node, error) {
	tree := make(map[string]*node)

	for _, r := range rs {
		if err := r.insert(tree); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func newTreeSelector(tree map[string]any) (treeSelector, error) {
	v := viper.New()

	v.SetDefault(prefRoute, tree[prefRoute])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", prefRoute, selector)
		} else {
			selector = prefRoute
		}

		val := v.Get(selector)

		switch v := val.(type) {

		case nil:
			return nil, fmt.Errorf("'%s' not found", selector)

		case map[string]any:
			return v, nil

		case string:
			return map[string]any{
				"path": v,
			}, nil

		default:
			return nil, fmt.Errorf(
				"invalid type for selector '%s': %T",
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
