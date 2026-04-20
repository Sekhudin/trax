package routes

import (
	"trax/internal/fs"
)

type routereader func(cfg *Config) ([]rawroute, error)

type builder struct{}

type generator struct {
	writer fs.FileWriter
}

type finalroute struct {
	raw      []rawroute
	routes   []route
	tree     map[string]*node
	selector treeselector
}

var (
	r   = route{}
	t   = tree{}
	bd  = builder{}
	gen = generator{
		writer: fs.NewOSWriter(),
	}
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

	tpl := newTemplate(&route.routes, route.selector, cfg)
	tp, err := tpl.build()
	if err != nil {
		return err
	}

	content := []byte(tp)

	return gen.generate(content, cfg)
}

func (g *generator) generate(content []byte, cfg *Config) error {
	return g.writer.Write(cfg.Output.Full, content)
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
