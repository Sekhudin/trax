package routes

import (
	"fmt"

	"github.com/spf13/viper"
)

type RoutesFile struct {
	Routes []Raw `mapstructure:"routes"`
}

func loadRoutesFile(c *RoutesConfig) ([]Raw, error) {
	v := viper.New()

	v.SetConfigFile(c.File.Full)

	v.SetDefault("foo", "")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read routes file: %w", err)
	}

	var rFile RoutesFile

	if err := v.Unmarshal(&rFile); err != nil {
		return nil, fmt.Errorf("failed to parse routes schema: %w", err)
	}

	if len(rFile.Routes) == 0 {
		return nil, fmt.Errorf("routes file is empty")
	}

	return rFile.Routes, nil
}

func ShowFromFile(c *RoutesConfig) (TreeSelector, error) {
	rFile, err := loadRoutesFile(c)
	if err != nil {
		return nil, err
	}

	routes, err := BuildRoutes(rFile)
	if err != nil {
		return nil, err
	}

	tree, err := BuildTree(routes)
	if err != nil {
		return nil, err
	}

	tSelector, err := NewTreeSelector(ToMap(tree))
	if err != nil {
		return nil, err
	}

	return tSelector, nil
}

func GenerateFromFile(c *RoutesConfig) (*[]Route, error) {
	rRoutes, err := loadRoutesFile(c)
	if err != nil {
		return nil, err
	}

	routes, err := BuildRoutes(rRoutes)
	if err != nil {
		return nil, err
	}

	return &routes, nil
}
