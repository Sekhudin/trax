package routes

import "fmt"

func loadFromDisc(c *RoutesConfig) ([]raw, error) {
	var rw []raw

	switch c.Strategy {
	case "next-page":
		w := walker{cfg: c, rule: &nextPageRule, wRule: stgNextPage{}}
		rs, err := w.walk()
		if err != nil {
			return nil, err
		}
		rw = rs

	case "next-app":
		w := walker{cfg: c, rule: &nextAppRule, wRule: stgNextApp{}}
		rs, err := w.walk()
		if err != nil {
			return nil, err
		}
		rw = rs

	default:
		return nil, fmt.Errorf("failed to read routes (strategy: %q)", c.Strategy)
	}

	return rw, nil
}

func ShowFromDisc(c *RoutesConfig) (TreeSelector, error) {
	rFile, err := loadFromDisc(c)
	if err != nil {
		return nil, err
	}

	routes, err := buildRoutes(rFile)
	if err != nil {
		return nil, err
	}

	tree, err := buildTree(routes)
	if err != nil {
		return nil, err
	}

	tSelector, err := newTreeSelector(toMap(tree))
	if err != nil {
		return nil, err
	}

	return tSelector, nil
}

func GenerateFromDisc(c *RoutesConfig) (*[]route, error) {
	rRoutes, err := loadFromDisc(c)
	if err != nil {
		return nil, err
	}

	routes, err := buildRoutes(rRoutes)
	if err != nil {
		return nil, err
	}

	return &routes, nil
}
