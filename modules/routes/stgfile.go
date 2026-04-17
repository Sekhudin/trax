package routes

import (
	"fmt"

	"github.com/spf13/viper"
)

type RoutesFile struct {
	Routes []RawRoute `mapstructure:"routes"`
}

func loadRoutesFile(c *RoutesConfig) ([]RawRoute, error) {
	v := viper.New()

	v.SetConfigFile(c.File.Full)

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

func ShowFromFile(c *RoutesConfig) (*map[string]any, error) {
	rRoutes, err := loadRoutesFile(c)
	if err != nil {
		return nil, err
	}

	routes, err := BuildRoutes(rRoutes)
	if err != nil {
		return nil, err
	}

	rTree, err := BuildRouteTree(routes)
	if err != nil {
		return nil, err
	}

	mapRTree := ToMap(rTree)

	return &mapRTree, nil
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
