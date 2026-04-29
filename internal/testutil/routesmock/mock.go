package routesmock

import (
	"errors"

	"github.com/sekhudin/trax/modules/routes"
)

type RoutesConfig struct {
	LoadCalled            bool
	IsFileStrategyCalled  bool
	IsValidStrategyCalled bool

	LoadFn            func() (*routes.Config, error)
	IsFileStrategyFn  func() bool
	IsValidStrategyFn func() bool
}

func (r *RoutesConfig) Reset() {
	r.LoadCalled = false
	r.IsFileStrategyCalled = false
	r.IsValidStrategyCalled = false

	r.LoadFn = func() (*routes.Config, error) {
		return &routes.Config{}, nil
	}

	r.IsFileStrategyFn = func() bool {
		return false
	}

	r.IsValidStrategyFn = func() bool {
		return false
	}
}

func (r *RoutesConfig) Load() (*routes.Config, error) {
	r.LoadCalled = true
	if r.LoadFn != nil {
		return r.LoadFn()
	}
	return &routes.Config{}, nil
}

func (r *RoutesConfig) IsFileStrategy() bool {
	r.IsFileStrategyCalled = true
	if r.IsFileStrategyFn != nil {
		return r.IsFileStrategyFn()
	}
	return false
}

func (r *RoutesConfig) IsValidStrategy() bool {
	r.IsValidStrategyCalled = true
	if r.IsValidStrategyFn != nil {
		return r.IsValidStrategyFn()
	}
	return false
}

type RawBuilder struct {
	BuildCalled bool
	BuildFn     func() ([]routes.RawRoute, error)
}

func (r *RawBuilder) Reset() {
	r.BuildCalled = false

	r.BuildFn = func() ([]routes.RawRoute, error) {
		return []routes.RawRoute{}, nil
	}
}

func (r *RawBuilder) Build() ([]routes.RawRoute, error) {
	r.BuildCalled = true
	if r.BuildFn != nil {
		return r.BuildFn()
	}
	return []routes.RawRoute{}, nil
}

type RouteBuilder struct {
	BuildCalled bool
	BuildFn     func([]routes.RawRoute) ([]routes.Route, error)
}

func (r *RouteBuilder) Reset() {
	r.BuildCalled = false

	r.BuildFn = func([]routes.RawRoute) ([]routes.Route, error) {
		return []routes.Route{}, nil
	}
}

func (r *RouteBuilder) Build(rws []routes.RawRoute) ([]routes.Route, error) {
	r.BuildCalled = true
	if r.BuildFn != nil {
		return r.BuildFn(rws)
	}
	return []routes.Route{}, nil
}

type TreeBuilder struct {
	BuildCalled bool
	BuildFn     func([]routes.Route) (map[string]*routes.Node, routes.TreeSelector, error)
}

func (t *TreeBuilder) Reset() {
	t.BuildCalled = false

	t.BuildFn = func([]routes.Route) (map[string]*routes.Node, routes.TreeSelector, error) {
		return map[string]*routes.Node{}, func(selector string) (map[string]any, error) {
			return map[string]any{}, nil
		}, nil
	}
}

func (t *TreeBuilder) Build(rs []routes.Route) (map[string]*routes.Node, routes.TreeSelector, error) {
	t.BuildCalled = true
	if t.BuildFn != nil {
		return t.BuildFn(rs)
	}

	return map[string]*routes.Node{}, func(selector string) (map[string]any, error) {
		return map[string]any{}, nil
	}, nil
}

type TemplateBuilder struct {
	BuildCalled bool
	BuildFn     func() (string, error)
}

func (t *TemplateBuilder) Reset() {
	t.BuildCalled = false

	t.BuildFn = func() (string, error) {
		return "content", nil
	}
}

func (t *TemplateBuilder) Build() (string, error) {
	t.BuildCalled = true
	if t.BuildFn != nil {
		return t.BuildFn()
	}
	return "content", nil
}

type Generator struct {
	GenerateCalled bool
	GenerateFn     func(string) error
}

func (g *Generator) Reset() {
	g.GenerateCalled = false

	g.GenerateFn = func(path string) error {
		return nil
	}
}

func (g *Generator) Generate(path string) error {
	g.GenerateCalled = true
	if g.GenerateFn != nil {
		return g.GenerateFn(path)
	}
	return nil
}

type Builder struct {
	BuildCalled bool
	BuildFn     func() (routes.BuildResult, error)
}

func (b *Builder) Reset() {
	b.BuildCalled = false

	b.BuildFn = func() (routes.BuildResult, error) {
		return &BuildResult{}, nil
	}
}

func (b *Builder) Build() (routes.BuildResult, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return &BuildResult{}, nil
}

type BuildResult struct {
	RawCalled      bool
	RoutesCalled   bool
	TreeCalled     bool
	SelectorCalled bool
	SelectCalled   bool

	RawFn      func() []routes.RawRoute
	RoutesFn   func() []routes.Route
	TreeFn     func() map[string]*routes.Node
	SelectorFn func() routes.TreeSelector
	SelectFn   func(key string) (map[string]any, error)
}

func (b *BuildResult) Reset() {
	b.RawCalled = false
	b.RoutesCalled = false
	b.TreeCalled = false
	b.SelectorCalled = false
	b.SelectCalled = false

	b.RawFn = func() []routes.RawRoute {
		return []routes.RawRoute{}
	}

	b.RoutesFn = func() []routes.Route {
		return []routes.Route{}
	}

	b.TreeFn = func() map[string]*routes.Node {
		return map[string]*routes.Node{}
	}

	b.SelectorFn = func() routes.TreeSelector {
		return func(selector string) (map[string]any, error) {
			return map[string]any{}, nil
		}
	}

	b.SelectFn = func(key string) (map[string]any, error) {
		if key == "fail" {
			return nil, errors.New("fail")
		}

		return map[string]any{}, nil
	}
}

func (b *BuildResult) Raw() []routes.RawRoute {
	b.RawCalled = true
	if b.RawFn != nil {
		return b.RawFn()
	}
	return []routes.RawRoute{}
}

func (b *BuildResult) Routes() []routes.Route {
	b.RoutesCalled = true
	if b.RoutesFn != nil {
		return b.RoutesFn()
	}
	return []routes.Route{}
}

func (b *BuildResult) Tree() map[string]*routes.Node {
	b.TreeCalled = true
	if b.TreeFn != nil {
		return b.TreeFn()
	}
	return map[string]*routes.Node{}
}

func (b *BuildResult) Selector() routes.TreeSelector {
	b.SelectorCalled = true
	if b.SelectorFn != nil {
		return b.SelectorFn()
	}

	return func(selector string) (map[string]any, error) {
		return map[string]any{}, nil
	}
}

func (b *BuildResult) Select(key string) (map[string]any, error) {
	b.SelectCalled = true
	if b.SelectFn != nil {
		return b.SelectFn(key)
	}

	if key == "fail" {
		return nil, errors.New("fail")
	}

	return map[string]any{}, nil
}
