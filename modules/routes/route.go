package routes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

type routerule struct {
	identPattern  *regexp.Regexp
	staticPattern *regexp.Regexp
}

type rawroute struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
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

type routefile struct {
	Routes []rawroute `mapstructure:"routes"`
}

var (
	nextApp  = nextapp{}
	nextPage = nextpage{}

	routeRule = routerule{
		identPattern:  regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`),
		staticPattern: regexp.MustCompile(`^[A-Za-z0-9_-]+$`),
	}
)

func (*route) readFile(c *Config) ([]rawroute, error) {
	v := viper.New()

	v.SetConfigFile(c.File.Full)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var rf routefile

	if err := v.Unmarshal(&rf); err != nil {
		return nil, err
	}

	if len(rf.Routes) == 0 {
		return nil, appErr.NewConfigNotFoundError("routes", "routes file is empty")
	}

	return rf.Routes, nil
}

func (*route) readDisc(cfg *Config) ([]rawroute, error) {
	switch cfg.Strategy {
	case "next-app":
		w := walker{strategy: &nextApp, config: cfg, rule: &nextRule.app}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	case "next-page":
		w := walker{strategy: &nextPage, config: cfg, rule: &nextRule.page}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	default:
		msg := fmt.Sprintf("failed to read routes (strategy: %q)", cfg.Strategy)
		return nil, appErr.NewValidationError("strategy", msg)
	}
}

func (*route) build(rw []rawroute) ([]route, error) {
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
		if !routeRule.identPattern.MatchString(cut) {
			msg := fmt.Sprintf("%q invalid param name: %s", r.name, cut)

			return "", "", appErr.NewValidationError("path", msg)
		}
		return "$" + cut, "param", nil
	}

	if !routeRule.staticPattern.MatchString(part) {
		msg := fmt.Sprintf("%q invalid path segment: %s", r.name, part)

		return "", "", appErr.NewValidationError("path", msg)
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
			msg := fmt.Sprintf("%q path wildcard route cannot coexist with other routes at the same level", r.name)

			return appErr.NewValidationError("path", msg)
		}

		if kind == "param" && c.kind == "static" {
			msg := fmt.Sprintf("%q path param route conflicts with existing static route", r.name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "static" && c.kind == "param" {
			msg := fmt.Sprintf("%q path static route conflicts with existing param route", r.name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "param" && c.kind == "param" {
			msg := fmt.Sprintf("%q path multiple param routes at the same level are not allowed", r.name)

			return appErr.NewValidationError("path", msg)
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
			msg := fmt.Sprintf("conflicting segment %q and %q produce same key", nd.segment, part)

			return appErr.NewValidationError("path", msg)
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
				msg := fmt.Sprintf("%q duplicate route detected: %q vs %q", r.name, nd.root, fPath)

				return appErr.NewValidationError("path", msg)

			}
			nd.root = fPath
		}

		current = nd.children
	}

	return nil
}

func (r *rawroute) cleanPath() (string, error) {
	r.Path = strings.TrimSpace(r.Path)

	if !strings.HasPrefix(r.Path, "/") {
		msg := fmt.Sprintf("%q path must start with %q", r.Name, "/")

		return "", appErr.NewValidationError("path", msg)
	}

	if strings.Contains(r.Path, "//") {
		msg := fmt.Sprintf("%q path contains double slash: %s", r.Name, r.Path)

		return "", appErr.NewValidationError("path", msg)
	}

	if i := strings.Index(r.Path, "?"); i > -1 {
		r.Path = r.Path[:i]
	}

	if len(r.Path) > 1 {
		r.Path = strings.TrimSuffix(r.Path, "/")
	}

	return r.Path, nil
}

func (r *rawroute) splitPath() []string {
	prefix := viper.GetString("routes.prefix")

	ps := strings.Split(fmt.Sprintf("%s%s", prefix, r.Path), "/")
	var result []string

	for _, p := range ps {
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func (r *rawroute) validateParts(parts []string) error {
	for i, p := range parts {
		if p == "*" && i != len(parts)-1 {
			msg := fmt.Sprintf("%q path wildcard must be last segment", r.Name)

			return appErr.NewValidationError("path", msg)
		}

		if strings.Contains(p, "*") && p != "*" {
			msg := fmt.Sprintf("%q path wildcard must be a single segment %q: %s", r.Name, p, "*")

			return appErr.NewValidationError("path", msg)
		}
	}
	return nil
}
