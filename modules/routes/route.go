package routes

import (
	"fmt"
	"strings"
)

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

func BuildRoutes(raw []Raw) ([]Route, error) {
	routes := make([]Route, 0, len(raw))

	for _, r := range raw {
		clean, err := r.cleanPath()
		if err != nil {
			return nil, err
		}

		parts := r.splitPath()
		if err := r.validateParts(parts); err != nil {
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

func (r *Route) normalizePart(part string) (string, string, error) {
	part = strings.ToLower(part)
	if part == "*" {
		return "$wildcard", "wildcard", nil
	}

	if cut, found := strings.CutPrefix(part, ":"); found {
		if !identRgx.MatchString(cut) {
			return "", "", fmt.Errorf("'%s' invalid param name: %s", r.Name, cut)
		}
		return "$" + cut, "param", nil
	}

	if !staticRgx.MatchString(part) {
		return "", "", fmt.Errorf("'%s' invalid path segment: %s", r.Name, part)
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

func (r *Route) validateChild(current map[string]*Node, kind string) error {
	for _, c := range current {
		if kind == "wildcard" || c.Kind == "wildcard" {
			return fmt.Errorf("'%s' path wildcard route cannot coexist with other routes at the same level", r.Name)
		}

		if kind == "param" && c.Kind == "static" {
			return fmt.Errorf("'%s' path param route conflicts with existing static route", r.Name)
		}

		if kind == "static" && c.Kind == "param" {
			return fmt.Errorf("'%s' path static route conflicts with existing param route", r.Name)
		}

		if kind == "param" && c.Kind == "param" {
			return fmt.Errorf("'%s' path multiple param routes at the same level are not allowed", r.Name)
		}
	}

	return nil
}

func (r *Route) insert(tree map[string]*Node) error {
	parts := r.Parts
	fullPath := r.Path
	current := tree

	for i, part := range parts {
		key, kind, err := r.normalizePart(part)
		if err != nil {
			return err
		}

		if err := r.validateChild(current, kind); err != nil {
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
