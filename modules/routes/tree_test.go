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

func TestNormalizePart_AllBranches(t *testing.T) {
	b := newTestTreeBuilder()

	r := Route{Name: "r"}

	k, kind, err := b.normalizePart(r, "*")
	if err != nil || kind != "wildcard" || k != "*" {
		t.Fatal("wildcard failed")
	}

	k, kind, err = b.normalizePart(r, ":id")
	if err != nil || kind != "param" {
		t.Fatal("param failed")
	}

	_, _, err = b.normalizePart(r, ":123")
	if err == nil {
		t.Fatal("expected invalid param")
	}

	_, _, err = b.normalizePart(r, "%%%")
	if err == nil {
		t.Fatal("expected invalid static")
	}

	k, kind, err = b.normalizePart(r, "user-profile")
	if err != nil || k != "userProfile" {
		t.Fatal("dash transform failed")
	}

	k, kind, err = b.normalizePart(r, "users")
	if err != nil || k != "users" || kind != "static" {
		t.Fatal("static failed")
	}
}

func TestValidateChild_AllConflicts(t *testing.T) {
	b := newTestTreeBuilder()
	r := Route{Name: "r"}

	current := map[string]*Node{
		"a": {Kind: "wildcard"},
	}
	if b.validateChild(r, current, "static") == nil {
		t.Fatal("wildcard conflict expected")
	}

	current = map[string]*Node{
		"a": {Kind: "static"},
	}
	if b.validateChild(r, current, "param") == nil {
		t.Fatal("param vs static conflict expected")
	}

	current = map[string]*Node{
		"a": {Kind: "param"},
	}
	if b.validateChild(r, current, "static") == nil {
		t.Fatal("static vs param conflict expected")
	}

	current = map[string]*Node{
		"a": {Kind: "param"},
	}
	if b.validateChild(r, current, "param") == nil {
		t.Fatal("param vs param conflict expected")
	}
}

func TestInsert_ConflictsAndDuplicate(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*Node{}

	r := Route{Name: "r1", Path: "/users", Parts: []string{"users"}}
	if err := b.insert(r, tr); err != nil {
		t.Fatal(err)
	}

	r2 := Route{Name: "r2", Path: "/users2", Parts: []string{"users"}}
	if err := b.insert(r2, tr); err == nil {
		t.Fatal("expected duplicate route error")
	}

	tr = map[string]*Node{}
	r3 := Route{Name: "r3", Path: "/user-profile", Parts: []string{"user-profile"}}
	r4 := Route{Name: "r4", Path: "/userProfile", Parts: []string{"userProfile"}}

	if err := b.insert(r3, tr); err != nil {
		t.Fatal(err)
	}
	if err := b.insert(r4, tr); err == nil {
		t.Fatal("expected conflicting key error")
	}
}

func TestToMap_WithRootAndChildren(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*Node{
		"users": {
			Root: "/users",
			Children: map[string]*Node{
				"list": {
					Root:     "/users/list",
					Children: map[string]*Node{},
				},
			},
		},
	}

	m := b.toMap(tr)

	if m["users"] == nil {
		t.Fatal("users missing")
	}
}

func TestCreateSelector_AllBranches(t *testing.T) {
	b := newTestTreeBuilder()
	b.cfg.Prefix = "api"

	tr := map[string]any{
		"api": map[string]any{
			"users": map[string]any{
				"path": "/users",
			},
			"plain": "/plain",
		},
	}

	sel := b.createSelector(tr)

	m, err := sel("users")
	if err != nil || m == nil {
		t.Fatal("map selector failed")
	}

	m, err = sel("plain")
	if err != nil || m["path"] != "/plain" {
		t.Fatal("string selector failed")
	}

	_, err = sel("notfound")
	if err == nil {
		t.Fatal("expected not found")
	}
}

func TestBuild_ErrorAndSuccess(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []Route{
		{Name: "a", Path: "/a", Parts: []string{"%%%"}},
	}

	_, _, err := b.Build(rs)
	if err == nil {
		t.Fatal("expected error")
	}

	rs = []Route{
		{Name: "a", Path: "/users", Parts: []string{"users"}},
	}

	_, sel, err := b.Build(rs)
	if err != nil || sel == nil {
		t.Fatal("build failed")
	}
}

func TestToMap_SkipNilNode(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*Node{
		"a": nil,
	}

	m := b.toMap(tr)

	if len(m) != 0 {
		t.Fatal("nil Node should be skipped")
	}
}

func TestInsert_ExistingKeyDifferentSegment(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*Node{}

	r1 := Route{Name: "r1", Path: "/x", Parts: []string{"user-s"}}
	if err := b.insert(r1, tr); err != nil {
		t.Fatal(err)
	}

	r2 := Route{Name: "r2", Path: "/x", Parts: []string{"userS"}}

	if err := b.insert(r2, tr); err == nil {
		t.Fatal("expected conflicting segment error")
	}
}

func TestCreateSelector_RootAndWeirdSuffixes(t *testing.T) {
	b := newTestTreeBuilder()
	b.cfg.Prefix = "api"

	tr := map[string]any{
		"api": map[string]any{
			"users": map[string]any{
				"$": "/users",
			},
			"num": 123,
		},
	}

	sel := b.createSelector(tr)

	if _, err := sel(""); err != nil {
		t.Fatal(err)
	}

	if _, err := sel("users$"); err != nil {
		t.Fatal(err)
	}
	if _, err := sel("users."); err != nil {
		t.Fatal(err)
	}
	if _, err := sel("users?"); err != nil {
		t.Fatal(err)
	}
	if _, err := sel("users?.$"); err != nil {
		t.Fatal(err)
	}

	if _, err := sel("num"); err == nil {
		t.Fatal("expected invalid type error")
	}
}

func TestBuild_DuplicateRouteSameLevel(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []Route{
		{Name: "r1", Path: "/users", Parts: []string{"users"}},
		{Name: "r2", Path: "/admin", Parts: []string{"users"}},
	}

	_, _, err := b.Build(rs)
	if err == nil {
		t.Fatal("expected duplicate route error")
	}
}

func TestBuild_WildcardConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []Route{
		{Name: "r1", Path: "/x/*", Parts: []string{"*"}},
		{Name: "r2", Path: "/x/users", Parts: []string{"users"}},
	}

	_, _, err := b.Build(rs)
	if err == nil {
		t.Fatal("expected wildcard conflict error")
	}
}

func TestBuild_ParamStaticConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []Route{
		{Name: "r1", Path: "/users", Parts: []string{"users"}},
		{Name: "r2", Path: "/:id", Parts: []string{":id"}},
	}

	_, _, err := b.Build(rs)
	if err == nil {
		t.Fatal("expected param/static conflict error")
	}
}

func TestBuild_MultipleParamConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []Route{
		{Name: "r1", Path: "/:id", Parts: []string{":id"}},
		{Name: "r2", Path: "/:name", Parts: []string{":name"}},
	}

	_, _, err := b.Build(rs)
	if err == nil {
		t.Fatal("expected multiple param conflict error")
	}
}
