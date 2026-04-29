package routes

import (
	"testing"

	"github.com/sekhudin/trax/internal/config"
)

func newTestTreeBuilder() tree {
	return tree{
		cfg: &Config{
			Symbols: &config.RoutesSymbols{
				Param:    ":",
				Wildcard: "*",
				Root:     "$",
			},
		},
	}
}

func TestTreeBuilder_Success(t *testing.T) {
	b := newTestTreeBuilder()
	r := Route{Name: "r"}

	t.Run("normalize_wildcard_part", func(t *testing.T) {
		k, kind, err := b.normalizePart(r, "*")
		if err != nil || kind != "wildcard" || k != "*" {
			t.Fatal("fail")
		}
	})

	t.Run("normalize_param_part", func(t *testing.T) {
		_, kind, err := b.normalizePart(r, ":id")
		if err != nil || kind != "param" {
			t.Fatal("fail")
		}
	})

	t.Run("transform_dash_segment", func(t *testing.T) {
		k, _, err := b.normalizePart(r, "user-profile")
		if err != nil || k != "userProfile" {
			t.Fatal("fail")
		}
	})

	t.Run("normalize_static_part", func(t *testing.T) {
		k, kind, err := b.normalizePart(r, "users")
		if err != nil || k != "users" || kind != "static" {
			t.Fatal("fail")
		}
	})

	t.Run("tomap_nested_children", func(t *testing.T) {
		tr := map[string]*Node{
			"users": {
				Root: "/users",
				Children: map[string]*Node{
					"list": {Root: "/users/list", Children: map[string]*Node{}},
				},
			},
		}
		m := b.toMap(tr)
		if m["users"] == nil {
			t.Fatal("fail")
		}
	})

	t.Run("selector_map_value", func(t *testing.T) {
		b.cfg.Prefix = "api"
		tr := map[string]any{"api": map[string]any{"users": map[string]any{"path": "/u"}}}
		sel := b.createSelector(tr)
		m, err := sel("users")
		if err != nil || m == nil {
			t.Fatal("fail")
		}
	})

	t.Run("selector_string_value", func(t *testing.T) {
		b.cfg.Prefix = "api"
		tr := map[string]any{"api": map[string]any{"plain": "/p"}}
		sel := b.createSelector(tr)
		m, err := sel("plain")
		if err != nil || m["path"] != "/p" {
			t.Fatal("fail")
		}
	})

	t.Run("build_valid_routes", func(t *testing.T) {
		rs := []Route{{Name: "a", Path: "/u", Parts: []string{"u"}}}
		_, sel, err := b.Build(rs)
		if err != nil || sel == nil {
			t.Fatal("fail")
		}
	})

	t.Run("selector_magic_suffixes", func(t *testing.T) {
		b.cfg.Prefix = "api"
		tr := map[string]any{"api": map[string]any{"users": map[string]any{"$": "/u"}}}
		sel := b.createSelector(tr)
		suffixes := []string{"", "users$", "users.", "users?", "users?.$"}
		for _, s := range suffixes {
			if _, err := sel(s); err != nil {
				t.Fatal("fail")
			}
		}
	})
}

func TestTreeBuilder_Error(t *testing.T) {
	b := newTestTreeBuilder()
	r := Route{Name: "r"}

	t.Run("invalid_param_name", func(t *testing.T) {
		if _, _, err := b.normalizePart(r, ":123"); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("invalid_static_segment", func(t *testing.T) {
		if _, _, err := b.normalizePart(r, "%%%"); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("wildcard_level_conflict", func(t *testing.T) {
		curr := map[string]*Node{"a": {Kind: "wildcard"}}
		if b.validateChild(r, curr, "static") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("param_static_conflict", func(t *testing.T) {
		curr := map[string]*Node{"a": {Kind: "static"}}
		if b.validateChild(r, curr, "param") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("static_param_conflict", func(t *testing.T) {
		curr := map[string]*Node{"a": {Kind: "param"}}
		if b.validateChild(r, curr, "static") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("multiple_param_conflict", func(t *testing.T) {
		curr := map[string]*Node{"a": {Kind: "param"}}
		if b.validateChild(r, curr, "param") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("duplicate_route_path", func(t *testing.T) {
		tr := map[string]*Node{}
		b.insert(Route{Path: "/u1", Parts: []string{"u"}}, tr)
		if b.insert(Route{Path: "/u2", Parts: []string{"u"}}, tr) == nil {
			t.Fatal("fail")
		}
	})

	t.Run("conflicting_key_segments", func(t *testing.T) {
		tr := map[string]*Node{}
		b.insert(Route{Parts: []string{"user-s"}}, tr)
		if b.insert(Route{Parts: []string{"userS"}}, tr) == nil {
			t.Fatal("fail")
		}
	})

	t.Run("selector_not_found", func(t *testing.T) {
		b.cfg.Prefix = "api"
		sel := b.createSelector(map[string]any{"api": map[string]any{}})
		if _, err := sel("notfound"); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("build_invalid_parts", func(t *testing.T) {
		rs := []Route{{Parts: []string{"%%%"}}}
		if _, _, err := b.Build(rs); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("selector_invalid_type", func(t *testing.T) {
		b.cfg.Prefix = "api"
		sel := b.createSelector(map[string]any{"api": map[string]any{"num": 123}})
		if _, err := sel("num"); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("build_duplicate_level", func(t *testing.T) {
		rs := []Route{
			{Path: "/u", Parts: []string{"u"}},
			{Path: "/a", Parts: []string{"u"}},
		}
		if _, _, err := b.Build(rs); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("build_wildcard_conflict", func(t *testing.T) {
		rs := []Route{
			{Path: "/x/*", Parts: []string{"*"}},
			{Path: "/x/u", Parts: []string{"u"}},
		}
		if _, _, err := b.Build(rs); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("build_param_conflict", func(t *testing.T) {
		rs := []Route{
			{Path: "/u", Parts: []string{"u"}},
			{Path: "/:id", Parts: []string{":id"}},
		}
		if _, _, err := b.Build(rs); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("build_multi_param", func(t *testing.T) {
		rs := []Route{
			{Path: "/:id", Parts: []string{":id"}},
			{Path: "/:name", Parts: []string{":name"}},
		}
		if _, _, err := b.Build(rs); err == nil {
			t.Fatal("fail")
		}
	})
}

func TestTreeBuilder_Fallback(t *testing.T) {
	b := newTestTreeBuilder()

	t.Run("skip_nil_node", func(t *testing.T) {
		m := b.toMap(map[string]*Node{"a": nil})
		if len(m) != 0 {
			t.Fatal("fail")
		}
	})
}
