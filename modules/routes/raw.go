package routes

import (
	"fmt"
	"strings"
)

type Raw struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

func (r *Raw) cleanPath() (string, error) {
	r.Path = strings.TrimSpace(r.Path)

	if !strings.HasPrefix(r.Path, "/") {
		return "", fmt.Errorf("'%s' path must start with '/': %s", r.Name, r.Path)
	}

	if strings.Contains(r.Path, "//") {
		return "", fmt.Errorf("'%s' path contains double slash: %s", r.Name, r.Path)
	}

	if i := strings.Index(r.Path, "?"); i > -1 {
		r.Path = r.Path[:i]
	}

	if len(r.Path) > 1 {
		r.Path = strings.TrimSuffix(r.Path, "/")
	}

	return r.Path, nil
}

func (r *Raw) splitPath() []string {
	parts := strings.Split(fmt.Sprintf("%s%s", prefRoute, r.Path), "/")
	var result []string

	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func (r *Raw) validateParts(parts []string) error {
	for i, p := range parts {
		if p == "*" && i != len(parts)-1 {
			return fmt.Errorf("'%s' path wildcard must be last segment", r.Name)
		}

		if strings.Contains(p, "*") && p != "*" {
			return fmt.Errorf("'%s' path wildcard must be a single segment '*': %s", r.Name, p)
		}
	}
	return nil
}
