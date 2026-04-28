package routes

import (
	"fmt"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type RawRouteFile struct {
	Routes []RawRoute `mapstructure:"routes"`
}

type RawRoute struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type RawRouteBuilder interface {
	Build() ([]RawRoute, error)
}

type rawroute struct {
	cfg *Config
}

func NewRawRouteBuilder(cfg *Config) RawRouteBuilder {
	return &rawroute{cfg: cfg}
}

func (b *rawroute) Build() ([]RawRoute, error) {
	var reader func() ([]RawRoute, error)

	if b.cfg.IsFileStrategy {
		reader = b.readFile
	} else {
		reader = b.readDisc
	}

	rws, err := reader()
	if err != nil {
		return nil, err
	}

	return rws, nil
}

func (b *rawroute) readFile() ([]RawRoute, error) {
	v := viper.New()

	v.SetConfigFile(b.cfg.File.Full)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var rf RawRouteFile

	if err := v.Unmarshal(&rf); err != nil {
		return nil, err
	}

	if len(rf.Routes) == 0 {
		return nil, appErr.NewConfigNotFoundError("routes", "routes file is empty")
	}

	return rf.Routes, nil
}

func (b *rawroute) readDisc() ([]RawRoute, error) {
	nextApp := newNextApp()
	nextPage := newNextPage()

	switch b.cfg.Strategy {
	case "next-app":
		w := walker{strategy: nextApp, config: b.cfg, rule: &nextApp.rule.app}
		rws, err := w.walk()
		if err != nil {
			return nil, err
		}

		return rws, nil

	case "next-page":
		w := walker{strategy: nextPage, config: b.cfg, rule: &nextApp.rule.app}
		rws, err := w.walk()
		if err != nil {
			return nil, err
		}

		return rws, nil

	default:
		msg := fmt.Sprintf("failed to read routes (strategy: %q)", b.cfg.Strategy)
		return nil, appErr.NewValidationError("strategy", msg)
	}
}
