package routes

func ShowFromDisc(cfg *Cfg) (TreeSelector, error) {
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

func GenerateFromDisc(cfg *Cfg) error {
	return nil
}
