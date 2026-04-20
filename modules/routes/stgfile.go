package routes

import "fmt"

func ShowFromFile(cfg *Cfg) (TreeSelector, error) {
	r := route{}

	rw, err := r.readFile(cfg)
	if err != nil {
		return nil, err
	}

	rs, err := r.build(rw)
	if err != nil {
		return nil, err
	}

	t := tree{}

	tr, err := t.build(rs)
	if err != nil {
		return nil, err
	}

	ts, err := t.newSelector(t.toMap(tr))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func GenerateFromFile(cfg *Cfg) error {
	r := route{}

	rw, err := r.readFile(cfg)
	if err != nil {
		return err
	}

	rs, err := r.build(rw)
	if err != nil {
		return err
	}

	t := tree{}

	tr, err := t.build(rs)
	if err != nil {
		return err
	}

	ts, err := t.newSelector(t.toMap(tr))
	if err != nil {
		return err
	}

	tpl := newTpl(&rs, &ts, cfg)
	fmt.Println(tpl.build())

	return nil
}
