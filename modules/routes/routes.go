package routes

import (
	"fmt"
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/viper"
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
	Segment  string
	Kind     string
	Children map[string]*Node
}

var (
	identRgx  = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	staticRgx = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

func BuildRoutes(raw []RawRoute) ([]Route, error) {
	routes := make([]Route, 0, len(raw))

	for _, r := range raw {
		clean, err := cleanPath(r)
		if err != nil {
			return nil, err
		}

		parts := split(fmt.Sprintf("routes%s", clean))

		if err := validate(r.Name, parts); err != nil {
			return nil, err
		}

		routes = append(routes, Route{
			Name:  r.Name,
			Path:  clean,
			Parts: parts,
		})
	}

	return routes, nil
}

func BuildRouteTree(routes []Route) (map[string]*Node, error) {
	tree := make(map[string]*Node)

	for _, r := range routes {
		if err := insert(tree, &r); err != nil {
			return nil, err
		}
	}

	return tree, nil
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

type SelectorFunc func(selector string) (map[string]any, error)

func NewTreeSelector(tree map[string]any) (SelectorFunc, error) {
	v := viper.New()

	rKey := "routes"
	v.SetDefault(rKey, tree[rKey])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", rKey, selector)
		} else {
			selector = rKey
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

func cleanPath(r RawRoute) (string, error) {
	p := strings.TrimSpace(r.Path)

	if !strings.HasPrefix(p, "/") {
		return "", fmt.Errorf("'%s' path must start with '/': %s", r.Name, p)
	}

	if strings.Contains(p, "//") {
		return "", fmt.Errorf("'%s' path contains double slash: %s", r.Name, p)
	}

	if i := strings.Index(p, "?"); i > -1 {
		p = p[:i]
	}

	if len(p) > 1 {
		p = strings.TrimSuffix(p, "/")
	}

	return p, nil
}

func validate(name string, parts []string) error {
	for i, p := range parts {
		if p == "*" && i != len(parts)-1 {
			return fmt.Errorf("'%s' path wildcard must be last segment", name)
		}

		if strings.Contains(p, "*") && p != "*" {
			return fmt.Errorf("'%s' path wildcard must be a single segment '*': %s", name, p)
		}
	}
	return nil
}

func validateChild(current map[string]*Node, kind string, name string) error {
	for _, c := range current {
		if kind == "wildcard" || c.Kind == "wildcard" {
			return fmt.Errorf("'%s' path wildcard route cannot coexist with other routes at the same level", name)
		}

		if kind == "param" && c.Kind == "static" {
			return fmt.Errorf("'%s' path param route conflicts with existing static route", name)
		}

		if kind == "static" && c.Kind == "param" {
			return fmt.Errorf("'%s' path static route conflicts with existing param route", name)
		}

		if kind == "param" && c.Kind == "param" {
			return fmt.Errorf("'%s' path multiple param routes at the same level are not allowed", name)
		}
	}

	return nil
}

func split(path string) []string {
	parts := strings.Split(strings.TrimSpace(path), "/")

	var res []string
	for _, p := range parts {
		if p != "" {
			res = append(res, p)
		}
	}

	return res
}

func normalize(part string, name string) (string, string, error) {
	part = strings.ToLower(part)
	if part == "*" {
		return "$wildcard", "wildcard", nil
	}

	if cut, found := strings.CutPrefix(part, ":"); found {
		if !identRgx.MatchString(cut) {
			return "", "", fmt.Errorf("'%s' invalid param name: %s", name, cut)
		}
		return "$" + cut, "param", nil
	}

	if !staticRgx.MatchString(part) {
		return "", "", fmt.Errorf("'%s' invalid path segment: %s", name, part)
	}

	if strings.Contains(part, "-") {
		parts := strings.Split(part, "-")
		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, ""), "static", nil
	}

	return part, "static", nil
}

func insert(tree map[string]*Node, r *Route) error {
	parts := r.Parts
	fullPath := r.Path
	current := tree

	for i, part := range parts {
		key, kind, err := normalize(part, r.Name)
		if err != nil {
			return err
		}

		if err := validateChild(current, kind, r.Name); err != nil {
			return err
		}

		node, ok := current[key]
		if ok && node.Segment != part {
			return fmt.Errorf("conflicting segment '%s' and '%s' produce same key", node.Segment, part)
		}

		if !ok {
			node = &Node{
				Children: make(map[string]*Node),
				Kind:     kind,
				Segment:  part,
			}
			current[key] = node
		}

		if i == len(parts)-1 {
			if node.Root != "" && node.Root != fullPath {
				return fmt.Errorf("'%s' duplicate route detected: '%s' vs '%s'", r.Name, node.Root, fullPath)
			}
			node.Root = fullPath
		}

		current = node.Children
	}

	return nil
}
