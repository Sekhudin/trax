package routes

import (
	"encoding/json"
	"strings"

	appErr "trax/internal/errors"
	"trax/internal/ts"
)

type template struct {
	routes   *[]route
	selector treeselector
	config   *Config
}

func newTemplate(r *[]route, s treeselector, cfg *Config) template {
	return template{routes: r, selector: s, config: cfg}
}

func (t *template) build() (string, error) {
	var b strings.Builder

	if t.isTypescrpt() {
		b.WriteString(t.warning())
		b.WriteString("\n\n")

		b.WriteString(t.etSearchParams())
		b.WriteString(";\n\n")

		b.WriteString(t.etRoutePattern())
		b.WriteString(";\n\n")

		b.WriteString(t.etExactParams())
		b.WriteString(";\n\n")

		b.WriteString(t.etReplaceParams())
		b.WriteString(";\n\n")

		b.WriteString(t.etWithQuery())
		b.WriteString(";\n\n")

		b.WriteString(t.etRouteBuilder())
		b.WriteString(";\n\n")
	}

	tree, err := t.routeTree()
	if err != nil {
		return "", err
	}

	b.WriteString(tree)
	b.WriteString(";\n\n")

	return b.String(), nil
}

func (t *template) isTypescrpt() bool {
	return strings.HasSuffix(t.config.Output.Ext, ".ts")
}

func (t *template) routerPatterns() []string {
	r := make([]string, 0, len(*t.routes))

	for _, p := range *t.routes {
		r = append(r, p.path)
	}

	return r
}

func (t *template) jsonTree() (string, error) {
	data, err := t.selector("")
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", appErr.NewInternalError("tree", "failed to marshal data to json", err)
	}

	return string(jsonData), nil
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

func (*template) etExactParams() string {
	var b strings.Builder

	b.WriteString("type ExtractParams<Path extends string> = Path extends ")
	b.WriteString(ts.ToTypeLiteral("${string}:${infer Param}/${infer Rest}"))
	b.WriteString("\n")

	b.WriteString("? Param | ExtractParams<")
	b.WriteString(ts.ToTypeLiteral("/${Rest}"))
	b.WriteString(">")
	b.WriteString("\n")

	b.WriteString(": Path extends ")
	b.WriteString(ts.ToTypeLiteral("${string}:${infer Param}"))
	b.WriteString("\n")

	b.WriteString("? Param")
	b.WriteString("\n")

	b.WriteString(": never")

	return b.String()
}

func (*template) etReplaceParams() string {
	var b strings.Builder

	b.WriteString("type ReplaceParams<")
	b.WriteString("\n")

	b.WriteString("Path extends string,")
	b.WriteString("\n")

	b.WriteString("Params extends Record<string, string>,")
	b.WriteString("\n")

	b.WriteString("> = Path extends ")
	b.WriteString(ts.ToTypeLiteral("${infer Start}:${infer Key}/${infer Rest}"))
	b.WriteString("\n")

	b.WriteString("?")
	b.WriteString(ts.ToTypeLiteral("${Start}${Params[Key]}/${ReplaceParams<Rest, Params>}"))
	b.WriteString("\n")

	b.WriteString(": Path extends ")
	b.WriteString(ts.ToTypeLiteral("${infer Start}:${infer Key}"))
	b.WriteString("\n")

	b.WriteString("?")
	b.WriteString(ts.ToTypeLiteral("${Start}${Params[Key]}"))
	b.WriteString("\n")

	b.WriteString(": Path")

	return b.String()
}

func (*template) etWithQuery() string {
	var b strings.Builder

	b.WriteString("export type WithQuery<")
	b.WriteString("\n")

	b.WriteString("Path extends string,")
	b.WriteString("\n")

	b.WriteString("Query extends SearchParams | undefined,")
	b.WriteString("\n")

	b.WriteString("> = Query extends undefined ? Path :")
	b.WriteString(ts.ToTypeLiteral("${Path}?${string}"))

	return b.String()
}

func (*template) etRouteBuilder() string {
	var b strings.Builder

	b.WriteString("export type RouteBuilder<Path extends string> =")
	b.WriteString("\n")

	b.WriteString("ExtractParams<Path> extends never")
	b.WriteString("\n")

	b.WriteString("? <Q extends SearchParams | undefined = undefined>(searchParams?: Q) => WithQuery<Path, Q>")
	b.WriteString("\n")

	b.WriteString(": <")
	b.WriteString("\n")

	b.WriteString("P extends Record<ExtractParams<Path>, string>,")
	b.WriteString("\n")

	b.WriteString("Q extends SearchParams | undefined = undefined,")
	b.WriteString("\n")

	b.WriteString(">(")
	b.WriteString("\n")

	b.WriteString("params: P,")
	b.WriteString("\n")

	b.WriteString("searchParams?: Q,")
	b.WriteString("\n")

	b.WriteString(") => WithQuery<ReplaceParams<Path, P>, Q>")

	return b.String()
}

func (*template) etRoutesLike() string {
	var b strings.Builder

	b.WriteString("export type RoutesLike<Tree> = {")
	b.WriteString("\n")

	b.WriteString("[Key in keyof Tree]: Tree[Key] extends string")
	b.WriteString("\n")

	b.WriteString("? RouteBuilder<Tree[Key]>")
	b.WriteString("\n")

	b.WriteString(": RoutesLike<Tree[Key]>")
	b.WriteString("\n")

	return b.String()
}

func (t *template) routeTree() (string, error) {
	var b strings.Builder

	if t.isTypescrpt() {
		b.WriteString("export type RouteTree = typeof routeTree;")
		b.WriteString("\n\n")
	}

	tree, err := t.jsonTree()
	if err != nil {
		return "", err
	}

	b.WriteString("export const routeTree = ")
	b.WriteString(tree)

	if t.isTypescrpt() {
		b.WriteString(" as const")
	}

	return b.String(), nil
}
