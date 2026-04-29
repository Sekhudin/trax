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

func (m *mockRawBuilder) Reset() {
	m.BuildCalled = false

	m.BuildFn = func() ([]RawRoute, error) {
		return []RawRoute{}, nil
	}
}

func (m *mockRawBuilder) Build() ([]RawRoute, error) {
	m.BuildCalled = true
	if m.BuildFn != nil {
		return m.BuildFn()
	}
	return []RawRoute{}, nil
}

type mockRouteBuilder struct {
	BuildCalled bool
	BuildFn     func([]RawRoute) ([]Route, error)
}

func (m *mockRouteBuilder) Reset() {
	m.BuildCalled = false

	m.BuildFn = func([]RawRoute) ([]Route, error) {
		return []Route{}, nil
	}
}

func (m *mockRouteBuilder) Build(r []RawRoute) ([]Route, error) {
	m.BuildCalled = true
	if m.BuildFn != nil {
		return m.BuildFn(r)
	}
	return []Route{}, nil
}

type mockTreeBuilder struct {
	BuildCalled bool

	BuildFn func([]Route) (map[string]*Node, TreeSelector, error)
}

func (m *mockTreeBuilder) Reset() {
	m.BuildCalled = false

	m.BuildFn = func([]Route) (map[string]*Node, TreeSelector, error) {
		return map[string]*Node{}, func(selector string) (map[string]any, error) {
			return map[string]any{}, nil
		}, nil
	}
}

func (m *mockTreeBuilder) Build(r []Route) (map[string]*Node, TreeSelector, error) {
	m.BuildCalled = true
	if m.BuildFn != nil {
		return m.BuildFn(r)
	}

	return map[string]*Node{}, func(selector string) (map[string]any, error) {
		if selector == "fail" {
			return nil, errors.New("fail")
		}

		return map[string]any{}, nil
	}, nil
}

type mockTemplateBuilder struct {
	BuildCalled bool

	BuildFn func() (string, error)
}

func (m *mockTemplateBuilder) Reset() {
	m.BuildCalled = false

	m.BuildFn = func() (string, error) {
		return "content", nil
	}
}

func (m *mockTemplateBuilder) Build() (string, error) {
	m.BuildCalled = true
	if m.BuildFn != nil {
		return m.BuildFn()
	}
	return "content", nil
}

func TestBuilder_Success(t *testing.T) {
	mockRaw := mockRawBuilder{}
	mockRoute := mockRouteBuilder{}
	mockTree := mockTreeBuilder{}

	t.Run("complete_build_flow", func(t *testing.T) {
		b := &builder{
			Deps: builderdeps{
				Raw:   &mockRaw,
				Route: &mockRoute,
				Tree:  &mockTree,
			},
		}

		res, err := b.Build()
		if err != nil || res == nil {
			t.Fatal("build_failed")
		}

		if res.Raw() == nil || res.Routes() == nil || res.Tree() == nil {
			t.Fatal("result_component_nil")
		}

		if _, err := res.Select("valid"); err != nil {
			t.Fatal("select_should_work")
		}

		if _, err := res.Select("fail"); err == nil {
			t.Fatal("select_should_fail")
		}
	})
}

func TestBuilder_Error(t *testing.T) {
	mockRaw := mockRawBuilder{}
	mockRoute := mockRouteBuilder{}
	mockTree := mockTreeBuilder{}

	t.Run("tree_build_failed", func(t *testing.T) {
		mockTree.BuildFn = func(r []Route) (map[string]*Node, TreeSelector, error) {
			return nil, nil, errors.New("err")
		}

		b := &builder{Deps: builderdeps{
			Raw:   &mockRaw,
			Route: &mockRoute,
			Tree:  &mockTree,
		}}

		if _, err := b.Build(); err == nil {
			t.Fatal("should_catch_tree_err")
		}
	})

	t.Run("route_build_failed", func(t *testing.T) {
		mockRoute.BuildFn = func(rr []RawRoute) ([]Route, error) {
			return nil, errors.New("err")
		}

		b := &builder{Deps: builderdeps{
			Raw:   &mockRaw,
			Route: &mockRoute,
		}}

		if _, err := b.Build(); err == nil {
			t.Fatal("should_catch_route_err")
		}
	})

	t.Run("raw_build_failed", func(t *testing.T) {
		mockRaw.BuildFn = func() ([]RawRoute, error) {
			return nil, errors.New("error")
		}

		b := &builder{Deps: builderdeps{
			Raw: &mockRaw,
		}}

		if _, err := b.Build(); err == nil {
			t.Fatal("should_catch_raw_err")
		}
	})
}

func TestBuilder_Fallback(t *testing.T) {
	mockCfg := Config{
		Strategy: "next-app",
		Root:     ".",
		Output: &path.FilePath{
			Full: "/tmp/out.ts",
		},
	}

	t.Run("new_builder_interface", func(t *testing.T) {
		var _ Builder = NewBuilder(&mockCfg)
	})

	t.Run("selector_accessor_check", func(t *testing.T) {
		mockResult := &buildresult{
			Result: &Routes{
				Raw:    []RawRoute{},
				Routes: []Route{},
				Tree:   map[string]*Node{},
				Selector: func(selector string) (map[string]any, error) {
					return map[string]any{}, nil
				},
			},
		}

		if mockResult.Selector() == nil {
			t.Fatal("accessor_failed")
		}
	})
}
