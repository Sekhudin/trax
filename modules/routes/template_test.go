package routes

import (
	"errors"
	"strings"
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/path"
)

func testTemplate_Config(ext string, noDeps bool) *Config {
	return &Config{
		NoDeps: noDeps,
		Output: &path.FilePath{Ext: ext},
		Symbols: &config.RoutesSymbols{
			Param:    ":",
			Wildcard: "*",
			Root:     "$",
		},
	}
}

func testTemplate_SelectorOK() TreeSelector {
	return func(_ string) (map[string]any, error) {
		return map[string]any{
			"api": map[string]any{
				"$": "/api",
			},
		}, nil
	}
}

func TestTemplate_Success(t *testing.T) {
	t.Run("build_typescript_mode", func(t *testing.T) {
		tpl := NewTemplate(TemplateDeps{
			Cfg:      testTemplate_Config(".ts", false),
			Selector: testTemplate_SelectorOK(),
			Routes:   []Route{{Path: "/posts/:id"}},
		})

		out, err := tpl.Build()
		if err != nil || out == "" {
			t.Fatal("should_build_ts_successfully")
		}

		if !strings.Contains(out, "type RoutePattern") || !strings.Contains(out, "as const") {
			t.Fatal("missing_typescript_definitions")
		}
	})

	t.Run("build_js_nodeps_mode", func(t *testing.T) {
		tpl := NewTemplate(TemplateDeps{
			Cfg:      testTemplate_Config(".js", true),
			Selector: testTemplate_SelectorOK(),
			Routes:   []Route{{Path: "/users"}},
		})

		out, err := tpl.Build()
		if err != nil || out == "" {
			t.Fatal("should_build_js_successfully")
		}

		if strings.Contains(out, "import qs") {
			t.Fatal("should_not_contain_qs_import_in_nodeps_mode")
		}
	})

	t.Run("serialize_logic_branches", func(t *testing.T) {
		instance := NewTemplate(TemplateDeps{
			Cfg: testTemplate_Config(".ts", false),
		}).(*template)

		data := map[string]any{
			"$":    "/admin",
			"name": "admin_page",
		}

		out := instance.serilizeRoutes(data, "", "admin")

		if !strings.Contains(out, "createRoute(tree.admin.$)") {
			t.Fatal("root_symbol_should_trigger_createroute")
		}
		if !strings.Contains(out, `"admin_page"`) {
			t.Fatal("static_string_should_be_quoted")
		}
	})
}

func TestTemplate_Error(t *testing.T) {
	t.Run("selector_failure", func(t *testing.T) {
		selErr := func(_ string) (map[string]any, error) { return nil, errors.New("fail") }
		tpl := NewTemplate(TemplateDeps{
			Cfg:      testTemplate_Config(".ts", false),
			Selector: selErr,
		})

		_, err := tpl.Build()
		if err == nil {
			t.Fatal("should_catch_selector_error")
		}
	})

	t.Run("marshal_failure", func(t *testing.T) {
		badSel := func(_ string) (map[string]any, error) {
			return map[string]any{"bad": make(chan int)}, nil
		}
		instance := NewTemplate(TemplateDeps{
			Cfg:      testTemplate_Config(".ts", false),
			Selector: badSel,
		}).(*template)

		_, err := instance.rTreeJSON()
		if err == nil {
			t.Fatal("should_catch_json_marshal_error")
		}
	})
}

func TestTemplate_Fallback(t *testing.T) {
	t.Run("empty_routes_pattern", func(t *testing.T) {
		instance := NewTemplate(TemplateDeps{
			Cfg:    testTemplate_Config(".ts", false),
			Routes: []Route{},
		}).(*template)

		out := instance.tRoutePattern()
		if !strings.Contains(out, "type RoutePattern =") {
			t.Fatal("should_still_render_type_definition")
		}
	})

	t.Run("nodeps_query_string_logic", func(t *testing.T) {
		instance := NewTemplate(TemplateDeps{
			Cfg: testTemplate_Config(".js", true),
		}).(*template)

		out := instance.fToQueryString()
		if !strings.Contains(out, "new URLSearchParams()") {
			t.Fatal("should_use_native_api_when_nodeps_is_true")
		}
	})
}
