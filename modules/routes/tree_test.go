package routes

import (
	"testing"

	"github.com/sekhudin/trax/internal/config"
)

func newTestTreeBuilder() treebuilder {
	return treebuilder{
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

	r := route{name: "r"}

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
	r := route{name: "r"}

	current := map[string]*node{
		"a": {kind: "wildcard"},
	}
	if b.validateChild(r, current, "static") == nil {
		t.Fatal("wildcard conflict expected")
	}

	current = map[string]*node{
		"a": {kind: "static"},
	}
	if b.validateChild(r, current, "param") == nil {
		t.Fatal("param vs static conflict expected")
	}

	current = map[string]*node{
		"a": {kind: "param"},
	}
	if b.validateChild(r, current, "static") == nil {
		t.Fatal("static vs param conflict expected")
	}

	current = map[string]*node{
		"a": {kind: "param"},
	}
	if b.validateChild(r, current, "param") == nil {
		t.Fatal("param vs param conflict expected")
	}
}

func TestInsert_ConflictsAndDuplicate(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*node{}

	r := route{name: "r1", path: "/users", parts: []string{"users"}}
	if err := b.insert(r, tr); err != nil {
		t.Fatal(err)
	}

	r2 := route{name: "r2", path: "/users2", parts: []string{"users"}}
	if err := b.insert(r2, tr); err == nil {
		t.Fatal("expected duplicate route error")
	}

	tr = map[string]*node{}
	r3 := route{name: "r3", path: "/user-profile", parts: []string{"user-profile"}}
	r4 := route{name: "r4", path: "/userProfile", parts: []string{"userProfile"}}

	if err := b.insert(r3, tr); err != nil {
		t.Fatal(err)
	}
	if err := b.insert(r4, tr); err == nil {
		t.Fatal("expected conflicting key error")
	}
}

func TestToMap_WithRootAndChildren(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*node{
		"users": {
			root: "/users",
			children: map[string]*node{
				"list": {
					root:     "/users/list",
					children: map[string]*node{},
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

	rs := []route{
		{name: "a", path: "/a", parts: []string{"%%%"}},
	}

	_, _, err := b.build(rs)
	if err == nil {
		t.Fatal("expected error")
	}

	rs = []route{
		{name: "a", path: "/users", parts: []string{"users"}},
	}

	_, sel, err := b.build(rs)
	if err != nil || sel == nil {
		t.Fatal("build failed")
	}
}

func TestToMap_SkipNilNode(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*node{
		"a": nil,
	}

	m := b.toMap(tr)

	if len(m) != 0 {
		t.Fatal("nil node should be skipped")
	}
}

func TestInsert_ExistingKeyDifferentSegment(t *testing.T) {
	b := newTestTreeBuilder()

	tr := map[string]*node{}

	r1 := route{name: "r1", path: "/x", parts: []string{"user-s"}}
	if err := b.insert(r1, tr); err != nil {
		t.Fatal(err)
	}

	r2 := route{name: "r2", path: "/x", parts: []string{"userS"}}

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

	rs := []route{
		{name: "r1", path: "/users", parts: []string{"users"}},
		{name: "r2", path: "/admin", parts: []string{"users"}},
	}

	_, _, err := b.build(rs)
	if err == nil {
		t.Fatal("expected duplicate route error")
	}
}

func TestBuild_WildcardConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []route{
		{name: "r1", path: "/x/*", parts: []string{"*"}},
		{name: "r2", path: "/x/users", parts: []string{"users"}},
	}

	_, _, err := b.build(rs)
	if err == nil {
		t.Fatal("expected wildcard conflict error")
	}
}

func TestBuild_ParamStaticConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []route{
		{name: "r1", path: "/users", parts: []string{"users"}},
		{name: "r2", path: "/:id", parts: []string{":id"}},
	}

	_, _, err := b.build(rs)
	if err == nil {
		t.Fatal("expected param/static conflict error")
	}
}

func TestBuild_MultipleParamConflict(t *testing.T) {
	b := newTestTreeBuilder()

	rs := []route{
		{name: "r1", path: "/:id", parts: []string{":id"}},
		{name: "r2", path: "/:name", parts: []string{":name"}},
	}

	_, _, err := b.build(rs)
	if err == nil {
		t.Fatal("expected multiple param conflict error")
	}
}
