package routes

import "fmt"

type routereader func(cfg *Config) ([]rawroute, error)

type builder struct{}

type finalroute struct {
	raw      []rawroute
	routes   []route
	tree     map[string]*node
	selector treeselector
}

var (
	r  = route{}
	t  = tree{}
	bd = builder{}
)

func Show(cfg *Config) (treeselector, error) {
	route, err := bd.build(cfg)
	if err != nil {
		return nil, err
	}

	return route.selector, nil
}

func Generate(cfg *Config) error {
	route, err := bd.build(cfg)
	if err != nil {
		return err
	}

	tpl := newTemplate(&route.routes, &route.selector, cfg)
	tp, err := tpl.build()

	fmt.Println(tp)
	return nil
}

func (*builder) build(cfg *Config) (*finalroute, error) {
	var reader routereader

	if cfg.IsFileStrategy() {
		reader = r.readFile
	} else {
		reader = r.readDisc
	}

	rw, err := reader(cfg)
	if err != nil {
		return nil, err
	}

	rs, err := r.build(rw)
	if err != nil {
		return nil, err
	}

	tr, err := t.build(rs)
	if err != nil {
		return nil, err
	}

	ts, err := t.newSelector(t.toMap(tr))
	if err != nil {
		return nil, err
	}

	return &finalroute{
		raw:      rw,
		routes:   rs,
		tree:     tr,
		selector: ts,
	}, nil
}
