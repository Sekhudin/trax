package routesmock

import "github.com/sekhudin/trax/modules/routes"

type RoutesConfig struct {
	LoadCalled            bool
	IsFileStrategyCalled  bool
	IsValidStrategyCalled bool

	LoadFn            func() (*routes.Config, error)
	IsFileStrategyFn  func() bool
	IsValidStrategyFn func() bool
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

func (b *RawBuilder) Build() ([]routes.RawRoute, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return []routes.RawRoute{}, nil
}

type RouteBuilder struct {
	BuildCalled bool
	BuildFn     func() ([]routes.Route, error)
}

func (b *RouteBuilder) Build(r []routes.RawRoute) ([]routes.Route, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return []routes.Route{}, nil
}

type TreeBuilder struct {
	BuildCalled bool
	BuildFn     func() (map[string]*routes.Node, routes.TreeSelector, error)
}

func (b *TreeBuilder) Build(r []routes.Route) (map[string]*routes.Node, routes.TreeSelector, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}

	return map[string]*routes.Node{}, func(selector string) (map[string]any, error) {
		return map[string]any{}, nil
	}, nil
}

type Builder struct {
	BuildCalled bool
	BuildFn     func() (routes.BuildResult, error)
}

func (b *Builder) Build() (routes.BuildResult, error) {
	b.BuildCalled = true
	if b.BuildFn != nil {
		return b.BuildFn()
	}
	return &BuildResult{}, nil
}

type Generator struct {
	GenerateCalled bool
	GenerateFn     func(string) error
}

func (g *Generator) Generate(path string) error {
	g.GenerateCalled = true
	if g.GenerateFn != nil {
		return g.GenerateFn(path)
	}
	return nil
}

type Template struct {
	BuildCalled bool
	BuildFn     func() (string, error)
}

func (t *Template) Build() (string, error) {
	t.BuildCalled = true
	if t.BuildFn != nil {
		return t.BuildFn()
	}
	return "", nil
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

func (r *BuildResult) Raw() []routes.RawRoute {
	r.RawCalled = true
	if r.RawFn != nil {
		return r.RawFn()
	}
	return []routes.RawRoute{}
}

func (r *BuildResult) Routes() []routes.Route {
	r.RoutesCalled = true
	if r.RoutesFn != nil {
		return r.RoutesFn()
	}
	return []routes.Route{}
}

func (r *BuildResult) Tree() map[string]*routes.Node {
	r.TreeCalled = true
	if r.TreeFn != nil {
		return r.TreeFn()
	}
	return map[string]*routes.Node{}
}

func (r *BuildResult) Selector() routes.TreeSelector {
	r.SelectorCalled = true
	if r.SelectorFn != nil {
		return r.SelectorFn()
	}

	return func(selector string) (map[string]any, error) {
		return map[string]any{}, nil
	}
}

func (r *BuildResult) Select(key string) (map[string]any, error) {
	r.SelectCalled = true
	if r.SelectFn != nil {
		return r.SelectFn(key)
	}
	return map[string]any{}, nil
}
