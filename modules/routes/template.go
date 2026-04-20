package routes

import (
	"strings"

	"trax/internal/ts"
)

type template struct {
	routes   *[]route
	selector *treeselector
	config   *Config
}

func newTemplate(r *[]route, s *treeselector, cfg *Config) template {
	return template{routes: r, selector: s, config: cfg}
}

func (t *template) build() (string, error) {
	var b strings.Builder

	b.WriteString(t.warning())
	b.WriteString("\n\n")
	b.WriteString(t.etSearchParams())
	b.WriteString(";\n\n")
	b.WriteString(t.etRoutePattern())
	b.WriteString(";\n\n")

	return b.String(), nil
}

func (t *template) routerPatterns() []string {
	r := make([]string, 0, len(*t.routes))

	for _, p := range *t.routes {
		r = append(r, p.path)
	}

	return r
}

func (*template) warning() string {
	return "/** [trax] AUTO-GENERATED — DO NOT TOUCH OR DELETE */"
}

func (*template) importModule(ext string) string {
	return ext
}

func (*template) etSearchParams() string {
	return "export type SearchParams = Record<string, string | number | boolean | undefined> | URLSearchParams"
}

func (t *template) etRoutePattern() string {
	var b strings.Builder

	b.WriteString("export type RoutePattern = ")
	b.WriteString(ts.OrganizeStringUnion(t.routerPatterns()))

	return b.String()
}

func (t *template) etExactParams() string {
	var b strings.Builder

	return b.String()
}
