package routes

import (
	"fmt"
	"maps"
	"sort"

	"github.com/spf13/viper"
)

type TreeSelector func(selector string) (map[string]any, error)

func BuildTree(routes []Route) (map[string]*Node, error) {
	tree := make(map[string]*Node)

	for _, r := range routes {
		if err := r.insert(tree); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func NewTreeSelector(tree map[string]any) (TreeSelector, error) {
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

func ToMap(nodes map[string]*Node) map[string]any {
	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make(map[string]any, len(nodes))

	for _, k := range keys {
		n := nodes[k]

		if n == nil {
			continue
		}

		m := make(map[string]any)

		if n.Root != "" {
			m["root"] = n.Root
		}

		if len(n.Children) > 0 {
			maps.Copy(m, ToMap(n.Children))
		}

		result[k] = m
	}

	return result
}
