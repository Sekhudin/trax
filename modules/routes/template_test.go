package routes

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/path"
)

func testSymbols() *config.RoutesSymbols {
	return &config.RoutesSymbols{
		Param:    ":",
		Wildcard: "*",
		Root:     "$",
	}
}

func testConfigTS() *Config {
	return &Config{
		NoDeps: false,
		Output: &path.FilePath{
			Ext: ".ts",
		},
		Symbols: testSymbols(),
	}
}

func testConfigJSNoDeps() *Config {
	return &Config{
		NoDeps: true,
		Output: &path.FilePath{
			Ext: ".js",
		},
		Symbols: testSymbols(),
	}
}

func mockSelectorOK() treeselector {
	return func(_ string) (map[string]any, error) {
		return map[string]any{
			"api": map[string]any{
				"$": "/api",
				"users": map[string]any{
					"$": "/api/users",
				},
			},
		}, nil
	}
}

func mockSelectorError() treeselector {
	return func(_ string) (map[string]any, error) {
		return nil, errors.New("selector failed")
	}
}

func newTpl(cfg *Config, sel treeselector, routes []route) *template {
	return &template{
		deps: TemplateDeps{
			Cfg:      cfg,
			Selector: sel,
			Routes:   routes,
		},
	}
}

func TestNewTemplate(t *testing.T) {
	cfg := testConfigTS()

	deps := TemplateDeps{
		Cfg:      cfg,
		Routes:   []route{{path: "/test"}},
		Selector: mockSelectorOK(),
	}

	tpl := NewTemplate(deps)

	if tpl == nil {
		t.Fatal("expected template instance")
	}

	if _, ok := tpl.(*template); !ok {
		t.Fatal("expected *template type")
	}
}

func TestRoutesTpl_Build_TS_Mode(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), []route{
		{path: "/users"},
		{path: "/posts/:id"},
	})

	out, err := tpl.Build()
	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("expected output")
	}
}

func TestRoutesTpl_Build_JS_NoDeps(t *testing.T) {
	cfg := testConfigJSNoDeps()

	tpl := newTpl(cfg, mockSelectorOK(), []route{
		{path: "/users"},
	})

	out, err := tpl.Build()
	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("expected output")
	}
}

func TestRoutesTpl_Build_SelectorError(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorError(), []route{})

	out, err := tpl.Build()
	if err == nil {
		t.Fatal(err)
	}

	if out != "" {
		t.Fatal("expected tree json")
	}
}

func TestRoutesTpl_Warning(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	w := tpl.warning()

	if w == "" {
		t.Fatal("warning should not be empty")
	}
}

func TestRoutesTpl_IsTypeScript(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	if !tpl.isTypescrpt() {
		t.Fatal("expected typescript mode")
	}
}

func TestRoutesTpl_IsNoDeps(t *testing.T) {
	cfg := testConfigJSNoDeps()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	if !tpl.isNoDeps() {
		t.Fatal("expected no deps mode")
	}
}

func TestRoutesTpl_ImportDeps(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	if tpl.importDeps() == "" {
		t.Fatal("expected import qs")
	}
}

func TestRoutesTpl_ImportDeps_NoDeps(t *testing.T) {
	cfg := testConfigJSNoDeps()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	if tpl.importDeps() != "" {
		t.Fatal("expected empty import")
	}
}

func TestRoutesTpl_RoutePattern(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), []route{
		{path: "/a"},
		{path: "/b"},
	})

	out := tpl.tRoutePattern()

	if out == "" {
		t.Fatal("expected route pattern type")
	}
}

func TestRoutesTpl_ExactParams(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.tExactParams()

	if out == "" {
		t.Fatal("expected exact params type")
	}
}

func TestRoutesTpl_ReplaceParams(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.tReplaceParams()

	if out == "" {
		t.Fatal("expected replace params type")
	}
}

func TestRoutesTpl_WithQuery(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.tWithQuery()

	if out == "" {
		t.Fatal("expected with query type")
	}
}

func TestRoutesTpl_RoutesLike(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.tRoutesLike()

	if out == "" {
		t.Fatal("expected routes like type")
	}
}

func TestRoutesTpl_CleanPath(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.fCleanPath()

	if out == "" {
		t.Fatal("expected cleanPath function")
	}
}

func TestRoutesTpl_FillParams(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.fFillParams()

	if out == "" {
		t.Fatal("expected fillParams function")
	}
}

func TestRoutesTpl_ToQueryString(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.fToQueryString()

	if out == "" {
		t.Fatal("expected toQueryString function")
	}
}

func TestRoutesTpl_WithQueryFn(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.fWithQuery()

	if out == "" {
		t.Fatal("expected withQuery function")
	}
}

func TestRoutesTpl_CreateRoute(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out := tpl.fCreateRoute()

	if out == "" {
		t.Fatal("expected createRoute function")
	}
}

func TestRoutesTpl_RTreeJSON(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	out, err := tpl.rTreeJSON()
	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("expected tree json")
	}
}

func TestTemplate_rTreeJSON_MarshalError(t *testing.T) {
	cfg := testConfigTS()

	selector := func(string) (map[string]any, error) {
		return map[string]any{
			"bad": make(chan int),
		}, nil
	}

	tpl := newTpl(cfg, selector, nil)

	_, err := tpl.rTreeJSON()
	if err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestRoutesTpl_RoutesJSON(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), []route{
		{path: "/users"},
	})

	out, err := tpl.rRoutesJSON()
	if err != nil {
		t.Fatal("expected error nil")
	}

	if out == "" {
		t.Fatal("expected routes json")
	}
}

func TestTemplate_rRoutesJSON_AllBranches(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, func(string) (map[string]any, error) {
		return map[string]any{
			"api": map[string]any{
				"$":     "/api",
				"users": "/api/users",
			},
		}, nil
	}, nil)

	out, err := tpl.rRoutesJSON()
	if err != nil {
		t.Fatal("expected error is nil")
	}

	if out == "" {
		t.Fatal("expected output")
	}
}

func TestRoutesTpl_rRoutesJSON_SelectorError(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorError(), nil)

	out, err := tpl.rRoutesJSON()
	if err == nil {
		t.Fatal("expected error is nil")
	}

	if out != "" {
		t.Fatal("expected output is empty")
	}
}

func TestRoutesTpl_SerializeRoutes(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	data := map[string]any{
		"api": map[string]any{
			"$": "/api",
		},
	}

	out := tpl.serilizeRoutes(data, "", "")

	if out == "" {
		t.Fatal("expected serialized output")
	}
}

func TestTemplate_serializeRoutes_AllBranches(t *testing.T) {
	cfg := testConfigTS()

	tpl := newTpl(cfg, mockSelectorOK(), nil)

	data := map[string]any{
		"$": "/root",
		"users": map[string]any{
			"$":       "/users",
			"profile": "/users/profile",
		},
	}

	out := tpl.serilizeRoutes(data, "", "")

	if out == "" {
		t.Fatal("expected output")
	}
}
