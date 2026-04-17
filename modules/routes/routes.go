package routes

import (
	"maps"
	"strings"
)

type RawRoute struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type Route struct {
	Name  string
	Path  string
	Parts []string
}

type Node struct {
	Root     string
	Children map[string]*Node
}

func BuildRoutes(raw []RawRoute) []Route {
	routes := make([]Route, 0, len(raw))

	for _, r := range raw {
		routes = append(routes, Route{
			Name:  r.Name,
			Path:  r.Path,
			Parts: split(r.Path),
		})
	}

	return routes
}

func BuildRouteTree(routes []Route) map[string]any {
	root := make(map[string]*Node)

	for _, r := range routes {
		insert(root, r.Parts, r.Path)
	}

	return toMap(root)
}

func split(path string) []string {
	parts := strings.Split(path, "/")

	var res []string
	for _, p := range parts {
		if p != "" {
			res = append(res, p)
		}
	}

	return res
}

func normalize(part string) string {
	if cut, found := strings.CutPrefix(part, ":"); found {
		part = cut
	}

	if strings.Contains(part, "-") {
		parts := strings.Split(part, "-")

		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}

		return strings.Join(parts, "")
	}

	return part
}

func insert(tree map[string]*Node, parts []string, fullPath string) {
	current := tree

	for i, part := range parts {
		key := normalize(part)

		node, ok := current[key]
		if !ok {
			node = &Node{
				Children: make(map[string]*Node),
			}
			current[key] = node
		}

		if i == len(parts)-1 {
			node.Root = fullPath
		}

		current = node.Children
	}
}

func toMap(nodes map[string]*Node) map[string]any {
	result := make(map[string]any, len(nodes))

	for k, n := range nodes {
		if n == nil {
			continue
		}

		m := make(map[string]any)

		if n.Root != "" {
			m["root"] = n.Root
		}

		if len(n.Children) > 0 {
			maps.Copy(m, toMap(n.Children))
		}

		result[k] = m
	}

	return result
}
