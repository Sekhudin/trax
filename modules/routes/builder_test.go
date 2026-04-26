package routes

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/path"
)

type mockRaw struct {
	err error
}

func (m mockRaw) build() ([]rawroute, error) {
	return []rawroute{{}}, m.err
}

type mockRoute struct {
	err error
}

func (m mockRoute) build([]rawroute) ([]route, error) {
	return []route{{}}, m.err
}

type mockTree struct {
	err error
}

func (m mockTree) build([]route) (map[string]*node, treeselector, error) {
	if m.err != nil {
		return nil, nil, m.err
	}

	return map[string]*node{"x": {}}, func(string) (map[string]any, error) {
		return map[string]any{}, nil
	}, nil
}

func cfg() *Config {
	return &Config{
		Strategy: "next-app",
		Root:     ".",
		Output: &path.FilePath{
			Full: "/tmp/out.ts",
		},
	}
}

func TestBuilder_Build_Success(t *testing.T) {
	b := &builder{
		deps: builderdeps{
			raw:   mockRaw{},
			route: mockRoute{},
			tree:  mockTree{},
		},
	}

	res, err := b.Build()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if res == nil {
		t.Fatal("result is nil")
	}

	if res.Selector == nil {
		t.Fatal("selector is nil")
	}

	if len(res.Raw) == 0 {
		t.Fatal("raw is empty")
	}
}

func TestBuilder_Build_RawError(t *testing.T) {
	b := &builder{
		deps: builderdeps{
			raw:   mockRaw{err: errors.New("raw error")},
			route: mockRoute{},
			tree:  mockTree{},
		},
	}

	_, err := b.Build()
	if err == nil {
		t.Fatal("expected raw error")
	}
}

func TestBuilder_Build_RouteError(t *testing.T) {
	b := &builder{
		deps: builderdeps{
			raw:   mockRaw{},
			route: mockRoute{err: errors.New("route error")},
			tree:  mockTree{},
		},
	}

	_, err := b.Build()
	if err == nil {
		t.Fatal("expected route error")
	}
}

func TestBuilder_Build_TreeError(t *testing.T) {
	b := &builder{
		deps: builderdeps{
			raw:   mockRaw{},
			route: mockRoute{},
			tree:  mockTree{err: errors.New("tree error")},
		},
	}

	_, err := b.Build()
	if err == nil {
		t.Fatal("expected tree error")
	}
}

func TestNewBuilder_ReturnsInterface(t *testing.T) {
	b := NewBuilder(cfg())

	res, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil {
		t.Fatal("result nil")
	}
}
