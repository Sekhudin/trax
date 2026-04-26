package routes

type finalroute struct {
	Raw      []rawroute
	Routes   []route
	Tree     map[string]*node
	Selector treeselector
}

type builderdeps struct {
	raw   rawroutebuilderItf
	route routebuilderItf
	tree  treebuilderItf
}

type builder struct {
	deps builderdeps
}

type Builder interface {
	Build() (*finalroute, error)
}

func NewBuilder(cfg *Config) Builder {
	return &builder{
		deps: builderdeps{
			raw:   newRawRouteBuilder(cfg),
			route: newRouteBuilder(cfg),
			tree:  newTreeBuilder(cfg),
		},
	}
}

func (b *builder) Build() (*finalroute, error) {
	rws, err := b.deps.raw.build()
	if err != nil {
		return nil, err
	}

	rs, err := b.deps.route.build(rws)
	if err != nil {
		return nil, err
	}

	tr, trs, err := b.deps.tree.build(rs)
	if err != nil {
		return nil, err
	}

	return &finalroute{
		Raw:      rws,
		Routes:   rs,
		Tree:     tr,
		Selector: trs,
	}, nil
}
