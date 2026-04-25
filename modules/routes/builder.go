package routes

import (
	"github.com/sekhudin/trax/internal/fs"
)

type finalroute struct {
	raw      []rawroute
	routes   []route
	tree     map[string]*node
	selector treeselector
}

type generator struct {
	writer fs.FileWriter
}

type builder struct{}

var (
	bd  = newBuilder()
	gen = newGenerator(fs.NewOSWriter())
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

func newBuilder() *builder {
	return &builder{}
}

func newGenerator(writer fs.FileWriter) *generator {
	return &generator{writer}
}

func (g *generator) generate(content []byte, cfg *Config) error {
	return g.writer.Write(cfg.Output.Full, content)
}

func (*builder) build(cfg *Config) (*finalroute, error) {
	rwb := newRawRouteBuilder(cfg)
	rb := newRouteBuilder(cfg)
	tb := newTreeBuilder(cfg)

	rws, err := rwb.build()
	if err != nil {
		return nil, err
	}

	rs, err := rb.build(rws)
	if err != nil {
		return nil, err
	}

	tr, trs, err := tb.build(rs)
	if err != nil {
		return nil, err
	}

	return &finalroute{
		raw:      rws,
		routes:   rs,
		tree:     tr,
		selector: trs,
	}, nil
}
