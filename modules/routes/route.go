package routes

import (
	"fmt"
	"strings"
)

type route struct {
	name  string
	path  string
	parts []string
}

type node struct {
	root     string
	segment  string
	kind     string
	children map[string]*node
}

func buildRoutes(rw []raw) ([]route, error) {
	rs := make([]route, 0, len(rw))

	for _, r := range rw {
		clean, err := r.cleanPath()
		if err != nil {
			return nil, err
		}

		parts := r.splitPath()
		if err := r.validateParts(parts); err != nil {
			return nil, err
		}

		rs = append(rs, route{
			name:  r.Name,
			path:  clean,
			parts: parts,
		})
	}

	return rs, nil
}

func (r *route) normalizePart(part string) (string, string, error) {
	part = strings.ToLower(part)
	if part == "*" {
		return "$wildcard", "wildcard", nil
	}

	if cut, found := strings.CutPrefix(part, ":"); found {
		if !identRgx.MatchString(cut) {
			return "", "", fmt.Errorf("'%s' invalid param name: %s", r.name, cut)
		}
		return "$" + cut, "param", nil
	}

	if !staticRgx.MatchString(part) {
		return "", "", fmt.Errorf("'%s' invalid path segment: %s", r.name, part)
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

func (r *route) validateChild(current map[string]*node, kind string) error {
	for _, c := range current {
		if kind == "wildcard" || c.kind == "wildcard" {
			return fmt.Errorf("'%s' path wildcard route cannot coexist with other routes at the same level", r.name)
		}

		if kind == "param" && c.kind == "static" {
			return fmt.Errorf("'%s' path param route conflicts with existing static route", r.name)
		}

		if kind == "static" && c.kind == "param" {
			return fmt.Errorf("'%s' path static route conflicts with existing param route", r.name)
		}

		if kind == "param" && c.kind == "param" {
			return fmt.Errorf("'%s' path multiple param routes at the same level are not allowed", r.name)
		}
	}

	return nil
}

func (r *route) insert(tree map[string]*node) error {
	parts := r.parts
	fPath := r.path
	current := tree

	for i, part := range parts {
		key, kind, err := r.normalizePart(part)
		if err != nil {
			return err
		}

		if err := r.validateChild(current, kind); err != nil {
			return err
		}

		nd, ok := current[key]
		if ok && nd.segment != part {
			return fmt.Errorf("conflicting segment '%s' and '%s' produce same key", nd.segment, part)
		}

		if !ok {
			nd = &node{
				children: make(map[string]*node),
				kind:     kind,
				segment:  part,
			}
			current[key] = nd
		}

		if i == len(parts)-1 {
			if nd.root != "" && nd.root != fPath {
				return fmt.Errorf("'%s' duplicate route detected: '%s' vs '%s'", r.name, nd.root, fPath)
			}
			nd.root = fPath
		}

		current = nd.children
	}

	return nil
}
