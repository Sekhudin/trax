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

	var r routesFile

	if err := v.Unmarshal(&r); err != nil {
		return nil, fmt.Errorf("failed to parse routes schema: %w", err)
	}

	if len(r.Routes) == 0 {
		return nil, fmt.Errorf("routes file is empty")
	}

	return r.Routes, nil
}

func ShowFromFile(c *RoutesConfig) (TreeSelector, error) {
	r, err := loadFromFile(c)
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

func GenerateFromFile(c *RoutesConfig) error {
	r, err := loadFromFile(c)
	if err != nil {
		return err
	}

	rs, err := buildRoutes(r)
	if err != nil {
		return err
	}

	fmt.Println(rs)

	return nil
}
