package routes

import (
	"fmt"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"
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

type RouteBuilder interface {
	Build([]RawRoute) ([]Route, error)
}

type route struct {
	cfg *Config
}

func NewRouteBuilder(cfg *Config) RouteBuilder {
	return &route{cfg: cfg}
}

func (b *route) Build(rws []RawRoute) ([]Route, error) {
	rs := make([]Route, 0, len(rws))

	for _, rw := range rws {
		clean, err := b.cleanPath(rw)
		if err != nil {
			return nil, err
		}

		parts := b.splitPath(rw)
		if err := b.validateParts(rw, parts); err != nil {
			return nil, err
		}

		rs = append(rs, Route{
			Name:  rw.Name,
			Path:  clean,
			Parts: parts,
		})
	}

	return rs, nil
}

func (b *route) cleanPath(rw RawRoute) (string, error) {
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

func (b *route) splitPath(rw RawRoute) []string {
	ps := strings.Split(fmt.Sprintf("%s%s", b.cfg.Prefix, rw.Path), "/")
	var result []string

	for _, p := range ps {
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func (b *route) validateParts(rw RawRoute, parts []string) error {
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
