package routes

import (
	"fmt"

	"github.com/spf13/viper"
)

type routesFile struct {
	Routes []raw `mapstructure:"routes"`
}

func loadFromFile(c *RoutesConfig) ([]raw, error) {
	v := viper.New()

	v.SetConfigFile(c.File.Full)

	v.SetDefault("foo", "")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read routes file: %w", err)
	}

	var rFile routesFile

	if err := v.Unmarshal(&rFile); err != nil {
		return nil, fmt.Errorf("failed to parse routes schema: %w", err)
	}

	if len(rFile.Routes) == 0 {
		return nil, fmt.Errorf("routes file is empty")
	}

	return rFile.Routes, nil
}

func ShowFromFile(c *RoutesConfig) (TreeSelector, error) {
	rFile, err := loadFromFile(c)
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

func GenerateFromFile(c *RoutesConfig) (*[]route, error) {
	rRoutes, err := loadFromFile(c)
	if err != nil {
		return nil, err
	}

	routes, err := buildRoutes(rRoutes)
	if err != nil {
		return nil, err
	}

	return &routes, nil
}
