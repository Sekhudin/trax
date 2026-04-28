package routes

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/path"
)

type mockRawBuilder struct {
	BuildCalled bool
	BuildFn     func() ([]RawRoute, error)
}

func (b *mockRawBuilder) Build() ([]RawRoute, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return []RawRoute{}, nil
}

type mockRouteBuilder struct {
	BuildCalled bool
	BuildFn     func() ([]Route, error)
}

func (b *mockRouteBuilder) Build(r []RawRoute) ([]Route, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return []Route{}, nil
}

type mockTreeBuilder struct {
	BuildCalled bool
	BuildFn     func() (map[string]*Node, TreeSelector, error)
}

func (b *mockTreeBuilder) Build(r []Route) (map[string]*Node, TreeSelector, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}

	return map[string]*Node{}, func(selector string) (map[string]any, error) {
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
		Deps: builderdeps{
			Raw:   &mockRawBuilder{},
			Route: &mockRouteBuilder{},
			Tree: &mockTreeBuilder{
				BuildFn: func() (map[string]*Node, TreeSelector, error) {
					return map[string]*Node{}, func(selector string) (map[string]any, error) {
						if selector == "invalid" {
							return nil, errors.New("notfound")
						}
						return map[string]any{}, nil
					}, nil
				},
			},
		},
	}

	res, err := b.Build()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if res == nil {
		t.Fatal("result is nil")
	}

	if res.Raw() == nil {
		t.Fatal("raw is nil")
	}

	if res.Routes() == nil {
		t.Fatal("routes is nil")
	}

	if res.Tree() == nil {
		t.Fatal("tree is nil")
	}

	if res.Selector() == nil {
		t.Fatal("selector is nil")
	}

	if _, err := res.Select(""); err != nil {
		t.Fatal("select is nil")
	}

	if _, err := res.Select("invalid"); err == nil {
		t.Fatal("expected error")
	}
}

func TestBuilder_Build_RawError(t *testing.T) {
	b := &builder{
		Deps: builderdeps{
			Raw: &mockRawBuilder{
				BuildFn: func() ([]RawRoute, error) {
					return nil, errors.New("raw error")
				},
			},
			Route: &mockRouteBuilder{},
			Tree:  &mockTreeBuilder{},
		},
	}

	_, err := b.Build()
	if err == nil {
		t.Fatal("expected raw error")
	}
}

func TestBuilder_Build_RouteError(t *testing.T) {
	b := &builder{
		Deps: builderdeps{
			Raw: &mockRawBuilder{},
			Route: &mockRouteBuilder{
				BuildFn: func() ([]Route, error) {
					return nil, errors.New("route error")
				},
			},
			Tree: &mockTreeBuilder{},
		},
	}

	_, err := b.Build()
	if err == nil {
		t.Fatal("expected route error")
	}
}

func TestBuilder_Build_TreeError(t *testing.T) {
	b := &builder{
		Deps: builderdeps{
			Raw:   &mockRawBuilder{},
			Route: &mockRouteBuilder{},
			Tree: &mockTreeBuilder{
				BuildFn: func() (map[string]*Node, TreeSelector, error) {
					return nil, nil, errors.New("tree error")
				},
			},
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
