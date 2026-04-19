package routes

import "fmt"

func ShowFromDisc(cfg *Config) (TreeSelector, error) {
	r := route{}

	rw, err := r.readDisc(cfg)
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

func GenerateFromDisc(cfg *Config) error {
	r := route{}

	rw, err := r.readDisc(cfg)
	if err != nil {
		return nil
	}

	rs, err := r.build(rw)
	if err != nil {
		return err
	}

	fmt.Println(rs)

	return nil
}
