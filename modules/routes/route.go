package routes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type routefile struct {
	Routes []raw `mapstructure:"routes"`
}

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

var (
	identRgx  = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	staticRgx = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

func (*route) readFile(c *Config) ([]raw, error) {
	v := viper.New()

	v.SetConfigFile(c.File.Full)

	v.SetDefault("foo", "")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read routes file: %w", err)
	}

	var rf routefile

	if err := v.Unmarshal(&rf); err != nil {
		return nil, fmt.Errorf("failed to parse routes schema: %w", err)
	}

	if len(rf.Routes) == 0 {
		return nil, fmt.Errorf("routes file is empty")
	}

	return rf.Routes, nil
}

func (*route) readDisc(c *Config) ([]raw, error) {
	switch c.Strategy {
	case "next-page":
		w := walker{cfg: c, rule: &nextPageRule, wRule: stgNextPage{}}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	case "next-app":
		w := walker{cfg: c, rule: &nextAppRule, wRule: stgNextApp{}}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	default:
		return nil, fmt.Errorf("failed to read routes (strategy: %q)", c.Strategy)
	}
}

func (*route) build(rw []raw) ([]route, error) {
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
			return "", "", fmt.Errorf("%q invalid param name: %s", r.name, cut)
		}
		return "$" + cut, "param", nil
	}

	if !staticRgx.MatchString(part) {
		return "", "", fmt.Errorf("%q invalid path segment: %s", r.name, part)
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
			return fmt.Errorf("%q path wildcard route cannot coexist with other routes at the same level", r.name)
		}

		if kind == "param" && c.kind == "static" {
			return fmt.Errorf("%q path param route conflicts with existing static route", r.name)
		}

		if kind == "static" && c.kind == "param" {
			return fmt.Errorf("%q path static route conflicts with existing param route", r.name)
		}

		if kind == "param" && c.kind == "param" {
			return fmt.Errorf("%q path multiple param routes at the same level are not allowed", r.name)
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
			return fmt.Errorf("conflicting segment %q and %q produce same key", nd.segment, part)
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
				return fmt.Errorf("%q duplicate route detected: %q vs %q", r.name, nd.root, fPath)
			}
			nd.root = fPath
		}

		current = nd.children
	}

	return nil
}
