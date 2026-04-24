package routes

import (
	"encoding/json"
	"fmt"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"

	"github.com/sekhudin/trax/internal/ts"
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

	b.WriteString(t.warning())
	b.WriteString("\n\n")
	b.WriteString(t.importDeps())

	if t.isTypescrpt() {
		b.WriteString(t.tParams())
		b.WriteString("\n\n")

		b.WriteString(t.tSearchParams())
		b.WriteString("\n\n")

		b.WriteString(t.tRoutePattern())
		b.WriteString("\n\n")

		b.WriteString(t.tExactParams())
		b.WriteString("\n\n")

		b.WriteString(t.tReplaceParams())
		b.WriteString("\n\n")

		b.WriteString(t.tWithQuery())
		b.WriteString("\n\n")

		b.WriteString(t.tRouteBuilder())
		b.WriteString("\n\n")

		b.WriteString(t.tRoutesLike())
		b.WriteString("\n\n")
	}

	rTree, err := t.rTreeJSON()
	if err != nil {
		return "", err
	}

	rRoutes, err := t.rRoutesJSON()
	if err != nil {
		return "", err
	}

	b.WriteString(t.fCleanPath())
	b.WriteString("\n\n")

	b.WriteString(t.fFillParams())
	b.WriteString("\n\n")

	b.WriteString(t.fToQueryString())
	b.WriteString("\n\n")

	b.WriteString(t.fWithQuery())
	b.WriteString("\n\n")

	b.WriteString(t.fCreateRoute())
	b.WriteString("\n\n")

	b.WriteString(rTree)
	b.WriteString("\n\n")

	b.WriteString(rRoutes)
	b.WriteString("\n\n")

	return b.String(), nil
}

func (t *template) isTypescrpt() bool {
	return strings.HasSuffix(t.config.Output.Ext, ".ts")
}

func (t *template) isNoDeps() bool {
	return t.config.NoDeps
}

func (*template) warning() string {
	return "/** [trax] AUTO-GENERATED — DO NOT TOUCH OR DELETE */"
}

func (t *template) importDeps() string {
	var b strings.Builder

	if !t.isNoDeps() {
		b.WriteString(`import qs from "qs";`)
		b.WriteString("\n\n")
	}

	return b.String()
}

func (*template) tParams() string {
	return "type Params = Record<string, string | number>;"
}

func (*template) tSearchParams() string {
	return "type SearchParams = Record<string, string | number | boolean | undefined> | URLSearchParams;"
}

func (t *template) tRoutePattern() string {
	r := make([]string, 0, len(*t.routes))
	for _, p := range *t.routes {
		r = append(r, p.path)
	}

	var b strings.Builder

	b.WriteString("type RoutePattern = ")
	b.WriteString(ts.OrganizeStringUnion(r))
	b.WriteString(";")

	return b.String()
}

func (*template) tExactParams() string {
	var b strings.Builder

	b.WriteString("type ExtractParams<Path extends string> = Path extends ")
	b.WriteString(ts.ToTypeLiteral("${string}:${infer Param}/${infer Rest}"))

	b.WriteString(" ? Param | ExtractParams<")
	b.WriteString(ts.ToTypeLiteral("/${Rest}"))
	b.WriteString("> : Path extends ")
	b.WriteString(ts.ToTypeLiteral("${string}:${infer Param}"))
	b.WriteString(" ? Param : never;")

	return b.String()
}

func (*template) tReplaceParams() string {
	var b strings.Builder

	b.WriteString("type ReplaceParams<Path extends string, Params extends Record<string, string>> = ")

	b.WriteString("Path extends ")
	b.WriteString(ts.ToTypeLiteral("${infer Start}:${infer Key}/${infer Rest}"))
	b.WriteString(" ? ")
	b.WriteString(ts.ToTypeLiteral("${Start}${Params[Key]}/${ReplaceParams<Rest, Params>}"))
	b.WriteString(" : Path extends ")
	b.WriteString(ts.ToTypeLiteral("${infer Start}:${infer Key}"))
	b.WriteString(" ? ")
	b.WriteString(ts.ToTypeLiteral("${Start}${Params[Key]}"))
	b.WriteString(" : Path;")

	return b.String()
}

func (*template) tWithQuery() string {
	var b strings.Builder

	b.WriteString("type WithQuery<Path extends string, Query extends SearchParams | undefined> = ")
	b.WriteString("Query extends undefined ? Path :")
	b.WriteString(ts.ToTypeLiteral("${Path}?${string}"))
	b.WriteString(";")

	return b.String()
}

func (*template) tRouteBuilder() string {
	var b strings.Builder

	b.WriteString("type RouteBuilder<Path extends string> = Path extends ")
	b.WriteString(ts.ToTypeLiteral("${infer Base}/*"))
	b.WriteString(" ? ExtractParams<Base> extends never")
	b.WriteString(" ? <Q extends SearchParams | undefined>")
	b.WriteString("(wildcardPath: string, searchParams?: Q) => string")
	b.WriteString(" : <P extends Record<ExtractParams<Base>, string>,")
	b.WriteString(" Q extends SearchParams | undefined>")
	b.WriteString("(wildcardPath: string, params: P, searchParams?: Q) => string")
	b.WriteString(" : ExtractParams<Path> extends never")
	b.WriteString(" ? <Q extends SearchParams | undefined>")
	b.WriteString("(searchParams?: Q) => WithQuery<Path, Q>")
	b.WriteString(" : <P extends Record<ExtractParams<Path>, string>, Q extends SearchParams | undefined>")
	b.WriteString("(params: P, searchParams?: Q) => WithQuery<ReplaceParams<Path, P>, Q>;")

	return b.String()
}

func (*template) tRoutesLike() string {
	var b strings.Builder

	b.WriteString("type RoutesLike<Tree> = {")
	b.WriteString(" [Key in keyof Tree]: Tree[Key] extends string")
	b.WriteString(" ? RouteBuilder<Tree[Key]>")
	b.WriteString(" : RoutesLike<Tree[Key]>;")
	b.WriteString(" }")

	return b.String()
}

func (t *template) fCleanPath() string {
	var b strings.Builder

	b.WriteString("function cleanPath")
	if t.isTypescrpt() {
		b.WriteString("(path: string) {\n")
	} else {
		b.WriteString("(path) {\n")
	}

	b.WriteString(`const matcher = /^\//;`)
	b.WriteString("\n")
	b.WriteString(`return path.replace(matcher, "");`)
	b.WriteString("\n}")

	return b.String()
}

func (t *template) fFillParams() string {
	var b strings.Builder

	b.WriteString("function fillParams")
	if t.isTypescrpt() {
		b.WriteString("(pattern: string, params: Params): string {\n")
	} else {
		b.WriteString("(pattern, params) {\n")
	}

	b.WriteString(`let result = pattern;`)
	b.WriteString("\nfor (const key in params) {")
	b.WriteString("\nconst regexPattern = new RegExp(")
	b.WriteString(ts.ToTypeLiteral(`:${key}\\b`))
	b.WriteString(`, "g");`)
	b.WriteString("\nresult = result.replace(regexPattern, String(params[key]));")
	b.WriteString("\n}")
	b.WriteString("\nreturn result;")
	b.WriteString("\n}")

	return b.String()
}

func (t *template) fToQueryString() string {
	var b strings.Builder

	b.WriteString("function toQueryString")
	if t.isTypescrpt() {
		b.WriteString("(query?: SearchParams) {\n")
	} else {
		b.WriteString("(query = {}) {\n")
	}

	b.WriteString(`if (!query) return "";`)
	b.WriteString("\n\n")

	b.WriteString("if (query instanceof URLSearchParams) {")
	b.WriteString("\nconst queryString = query.toString();")
	b.WriteString("\nreturn queryString ? ")
	b.WriteString(ts.ToTypeLiteral("?${queryString}"))
	b.WriteString(` : "";`)
	b.WriteString("\n}")
	b.WriteString("\n\n")

	if t.isNoDeps() {
		b.WriteString("const params = new URLSearchParams();")
		b.WriteString("\nfor (const [key, value] of Object.entries(query)) {")
		b.WriteString("\nif (value !== undefined && value !== null) {")
		b.WriteString("\nparams.append(key, String(value));")
		b.WriteString("\n}")
		b.WriteString("\n}")
		b.WriteString("\n\n")
		b.WriteString("\nconst queryString = params.toString();")
		b.WriteString("return queryString ?")
		b.WriteString(ts.ToTypeLiteral("?${queryString}"))
		b.WriteString(` : "";`)
	} else {
		b.WriteString("\nconst queryString = qs.stringify(query, {")
		b.WriteString("\nskipNulls: true,")
		b.WriteString("\naddQueryPrefix: true,")
		b.WriteString("\n});")
		b.WriteString("\n\n")
		b.WriteString("return queryString;")
	}

	b.WriteString("\n}")

	return b.String()
}

func (t *template) fWithQuery() string {
	var b strings.Builder

	b.WriteString("function withQuery")
	if t.isTypescrpt() {
		b.WriteString("<Pattern extends string, Query extends SearchParams | undefined>")
		b.WriteString("(pattern: Pattern, query: Query)")
		b.WriteString(" : WithQuery<Pattern, Query> {")
		b.WriteString("\nif (!query) return pattern as WithQuery<Pattern, Query>;")
		b.WriteString("\n\n")
	} else {
		b.WriteString("(pattern, query) {")
		b.WriteString("\nif (!query) return ")
		b.WriteString(ts.ToTypeLiteral("${pattern}"))
		b.WriteString(";")
		b.WriteString("\n\n")
	}

	b.WriteString("const queryString = toQueryString(query);")
	b.WriteString("\n\n")

	b.WriteString("return ")
	b.WriteString(ts.ToTypeLiteral("${pattern}${queryString}"))

	if t.isTypescrpt() {
		b.WriteString(" as WithQuery<Pattern, Query>")
	}

	b.WriteString(";")
	b.WriteString("\n")

	b.WriteString("}")

	return b.String()
}

func (t *template) fCreateRoute() string {
	var b strings.Builder

	b.WriteString("function createRoute")
	if t.isTypescrpt() {
		b.WriteString("<Pattern extends RoutePattern>(pattern: Pattern): RouteBuilder<Pattern> {\n")
	} else {
		b.WriteString("(pattern) {\n")
	}

	b.WriteString(`const hasWildcard = pattern.endsWith("/*");`)
	b.WriteString("\nconst hasParams = /:[^/]+/.test(pattern);")
	b.WriteString("\n\n")

	b.WriteString("\nif (hasWildcard && hasParams) {")
	if t.isTypescrpt() {
		b.WriteString("return ((path: string, params: Params, searchParams?: SearchParams) => {")
	} else {
		b.WriteString("return (path, params, searchParams = {}) => {")
	}
	b.WriteString("\nconst clean = cleanPath(path);")
	b.WriteString(`const result = pattern.replace("/*",`)
	b.WriteString(ts.ToTypeLiteral("/${clean}"))
	b.WriteString(");\n")
	b.WriteString("\nreturn withQuery(fillParams(result, params), searchParams);")
	b.WriteString("\n}")
	if t.isTypescrpt() {
		b.WriteString(") as RouteBuilder<Pattern>")
	}
	b.WriteString(";\n}\n")

	b.WriteString("\nif (hasWildcard) {\n")
	if t.isTypescrpt() {
		b.WriteString("return ((path: string, searchParams?: SearchParams) => {")
	} else {
		b.WriteString("return (path, searchParams = {}) => {")
	}
	b.WriteString("\nconst clean = cleanPath(path);")
	b.WriteString(`const result = pattern.replace("/*",`)
	b.WriteString(ts.ToTypeLiteral("/${clean}"))
	b.WriteString(");\n")
	b.WriteString("\nreturn withQuery(result, searchParams);")
	b.WriteString("\n}")
	if t.isTypescrpt() {
		b.WriteString(") as RouteBuilder<Pattern>")
	}
	b.WriteString(";\n}\n")

	b.WriteString("\nif (hasParams) {\n")
	if t.isTypescrpt() {
		b.WriteString("return ((params: Params, searchParams?: SearchParams) => {")
	} else {
		b.WriteString("return (params, searchParams = {}) => {")
	}
	b.WriteString("\nreturn withQuery(fillParams(pattern, params), searchParams);")
	b.WriteString("\n}")
	if t.isTypescrpt() {
		b.WriteString(") as RouteBuilder<Pattern>")
	}
	b.WriteString(";\n}\n")

	if t.isTypescrpt() {
		b.WriteString("\nreturn ((searchParams?: SearchParams) => {")
	} else {
		b.WriteString("\nreturn (searchParams = {}) => {")
	}
	b.WriteString("\nreturn withQuery(pattern, searchParams);")
	b.WriteString("\n}")
	if t.isTypescrpt() {
		b.WriteString(") as RouteBuilder<Pattern>")
	}

	b.WriteString("\n")
	b.WriteString("}")

	return b.String()
}

func (t *template) rTreeJSON() (string, error) {
	data, err := t.selector("")
	if err != nil {
		return "", err
	}

	tree, err := json.Marshal(data)
	if err != nil {
		return "", appErr.NewInternalError("tree", "failed to marshal data to json", err)
	}

	var b strings.Builder

	b.WriteString("export const tree = ")
	b.WriteString(string(tree))

	if t.isTypescrpt() {
		b.WriteString(" as const")
	}

	b.WriteString(";")

	return b.String(), nil
}

func (t *template) serilizeRoutes(data map[string]any, indent string, currentPath string) string {
	var b strings.Builder
	b.WriteString("{\n")

	for key, value := range data {
		newPath := key
		if currentPath != "" {
			newPath = currentPath + "." + key
		}

		b.WriteString(indent + "  " + key + ": ")

		switch v := value.(type) {
		case string:
			if key == "root" {
				val := fmt.Sprintf("createRoute(tree.%s)", newPath)
				b.WriteString(val)
			} else {
				val := fmt.Sprintf("%q", v)
				b.WriteString(val)
			}
		case map[string]any:
			b.WriteString(t.serilizeRoutes(v, indent+"  ", newPath))
		}
		b.WriteString(",\n")
	}

	b.WriteString(indent + "}")
	return b.String()
}

func (t *template) rRoutesJSON() (string, error) {
	data, err := t.selector("")
	if err != nil {
		return "", err
	}

	var b strings.Builder
	routes := t.serilizeRoutes(data, "", "")

	if t.isTypescrpt() {
		b.WriteString("export const routes: RoutesLike<typeof tree> =")
	} else {
		b.WriteString("export const routes =")
	}

	b.WriteString(routes)
	b.WriteString(";")

	return b.String(), nil
}
