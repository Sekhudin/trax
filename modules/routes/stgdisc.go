package routes

import "fmt"

func loadFromDisc(c *RoutesConfig) ([]raw, error) {
	switch c.Strategy {
	case "next-page":
		w := walker{cfg: c, rule: &nextPageRule, wRule: stgNextPage{}}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	case "next-app":
		w := walker{cfg: c, rule: &nextAppRule, wRule: stgNextApp{}}
		r, err := w.walk()
		if err != nil {
			return nil, err
		}

		return r, nil

	default:
		return nil, fmt.Errorf("failed to read routes (strategy: %q)", c.Strategy)
	}
}

func ShowFromDisc(c *RoutesConfig) (TreeSelector, error) {
	r, err := loadFromDisc(c)
	if err != nil {
		return nil, err
	}

	rs, err := buildRoutes(r)
	if err != nil {
		return nil, err
	}

	tr, err := buildTree(rs)
	if err != nil {
		return nil, err
	}

	ts, err := newTreeSelector(toMap(tr))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func GenerateFromDisc(c *RoutesConfig) error {
	r, err := loadFromDisc(c)
	if err != nil {
		return nil
	}

	rs, err := buildRoutes(r)
	if err != nil {
		return nil
	}

	fmt.Println(rs)

	return nil
}
