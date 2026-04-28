package routes

type Builder interface {
	Build() (BuildResult, error)
}

type BuildResult interface {
	Raw() []RawRoute
	Routes() []Route
	Tree() map[string]*Node
	Selector() TreeSelector
	Select(key string) (map[string]any, error)
}

type Routes struct {
	Raw      []RawRoute
	Routes   []Route
	Tree     map[string]*Node
	Selector TreeSelector
}

type builderdeps struct {
	Raw   RawRouteBuilder
	Route RouteBuilder
	Tree  TreeBuilder
}

type builder struct {
	Deps builderdeps
}

type buildresult struct {
	Result *Routes
}

func NewBuilder(cfg *Config) Builder {
	return &builder{
		Deps: builderdeps{
			Raw:   NewRawRouteBuilder(cfg),
			Route: NewRouteBuilder(cfg),
			Tree:  NewTreeBuilder(cfg),
		},
	}
}

func (b *builder) Build() (BuildResult, error) {
	rws, err := b.Deps.Raw.Build()
	if err != nil {
		return nil, err
	}

	rs, err := b.Deps.Route.Build(rws)
	if err != nil {
		return nil, err
	}

	tr, trs, err := b.Deps.Tree.Build(rs)
	if err != nil {
		return nil, err
	}

	return &buildresult{
		Result: &Routes{
			Raw:      rws,
			Routes:   rs,
			Tree:     tr,
			Selector: trs,
		},
	}, nil
}

func (r *buildresult) Raw() []RawRoute {
	return r.Result.Raw
}

func (r *buildresult) Routes() []Route {
	return r.Result.Routes
}

func (r *buildresult) Tree() map[string]*Node {
	return r.Result.Tree
}

func (r *buildresult) Selector() TreeSelector {
	return r.Result.Selector
}

func (r *buildresult) Select(key string) (map[string]any, error) {
	res, err := r.Result.Selector(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
