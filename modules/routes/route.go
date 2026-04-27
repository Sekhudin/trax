package routes

import (
	"fmt"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"
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

type routebuilderItf interface {
	build([]rawroute) ([]route, error)
}

type routebuilder struct {
	cfg *Config
}

func newRouteBuilder(cfg *Config) routebuilderItf {
	return &routebuilder{cfg: cfg}
}

func (b *routebuilder) build(rws []rawroute) ([]route, error) {
	rs := make([]route, 0, len(rws))

	for _, rw := range rws {
		clean, err := b.cleanPath(rw)
		if err != nil {
			return nil, err
		}

		parts := b.splitPath(rw)
		if err := b.validateParts(rw, parts); err != nil {
			return nil, err
		}

		rs = append(rs, route{
			name:  rw.Name,
			path:  clean,
			parts: parts,
		})
	}

	return rs, nil
}

func (b *routebuilder) cleanPath(rw rawroute) (string, error) {
	p := strings.TrimSpace(rw.Path)

	if !strings.HasPrefix(p, "/") {
		msg := fmt.Sprintf("%q path must start with %q", rw.Name, "/")

		return "", appErr.NewValidationError("path", msg)
	}

	if strings.Contains(p, "//") {
		msg := fmt.Sprintf("%q path contains double slash: %s", rw.Name, p)

		return "", appErr.NewValidationError("path", msg)
	}

	if i := strings.Index(p, "?"); i > -1 {
		p = p[:i]
	}

	if len(p) > 1 {
		p = strings.TrimSuffix(p, "/")
	}

	return p, nil
}

func (b *routebuilder) splitPath(rw rawroute) []string {
	ps := strings.Split(fmt.Sprintf("%s%s", b.cfg.Prefix, rw.Path), "/")
	var result []string

	for _, p := range ps {
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func (b *routebuilder) validateParts(rw rawroute, parts []string) error {
	for i, p := range parts {
		if p == "*" && i != len(parts)-1 {
			msg := fmt.Sprintf("%q path wildcard must be last segment", rw.Name)

			return appErr.NewValidationError("path", msg)
		}

		if strings.Contains(p, "*") && p != "*" {
			msg := fmt.Sprintf("%q path wildcard must be a single segment %q: %s", rw.Name, p, "*")

			return appErr.NewValidationError("path", msg)
		}
	}
	return nil
}
