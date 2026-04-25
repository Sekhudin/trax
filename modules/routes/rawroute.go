package routes

import (
	"fmt"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type routefile struct {
	Routes []rawroute `mapstructure:"routes"`
}

type rawroute struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type rawroutebuilder struct {
	cfg *Config
}

var (
	nextApp  = nextapp{}
	nextPage = nextpage{}
)

func newRawRouteBuilder(cfg *Config) *rawroutebuilder {
	return &rawroutebuilder{cfg: cfg}
}

func (b *rawroutebuilder) build() ([]rawroute, error) {
	var reader func() ([]rawroute, error)

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

func (b *rawroutebuilder) readFile() ([]rawroute, error) {
	v := viper.New()

	v.SetConfigFile(b.cfg.File.Full)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var rf routefile

	if err := v.Unmarshal(&rf); err != nil {
		return nil, err
	}

	if len(rf.Routes) == 0 {
		return nil, appErr.NewConfigNotFoundError("routes", "routes file is empty")
	}

	return rf.Routes, nil
}

func (b *rawroutebuilder) readDisc() ([]rawroute, error) {
	switch b.cfg.Strategy {
	case "next-app":
		w := walker{strategy: &nextApp, config: b.cfg, rule: &nextRule.app}
		rws, err := w.walk()
		if err != nil {
			return nil, err
		}

		return rws, nil

	case "next-page":
		w := walker{strategy: &nextPage, config: b.cfg, rule: &nextRule.page}
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
