package routes

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type raw struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

func (r *raw) cleanPath() (string, error) {
	r.Path = strings.TrimSpace(r.Path)

	if !strings.HasPrefix(r.Path, "/") {
		return "", fmt.Errorf("%q path must start with %q: %s", r.Name, r.Path, "/")
	}

	if strings.Contains(r.Path, "//") {
		return "", fmt.Errorf("%q path contains double slash: %s", r.Name, r.Path)
	}

	if i := strings.Index(r.Path, "?"); i > -1 {
		r.Path = r.Path[:i]
	}

	if len(r.Path) > 1 {
		r.Path = strings.TrimSuffix(r.Path, "/")
	}

	return r.Path, nil
}

func (r *raw) splitPath() []string {
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

func (r *raw) validateParts(parts []string) error {
	for i, p := range parts {
		if p == "*" && i != len(parts)-1 {
			return fmt.Errorf("%q path wildcard must be last segment", r.Name)
		}

		if strings.Contains(p, "*") && p != "*" {
			return fmt.Errorf("%q path wildcard must be a single segment %q: %s", r.Name, p, "*")
		}
	}
	return nil
}
