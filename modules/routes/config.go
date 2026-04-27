package routes

import (
	"fmt"
	"path/filepath"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/path"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type RoutesConfig interface {
	Load() (*Config, error)
	IsFileStrategy() bool
	IsValidStartegy() bool
}

type Config struct {
	Strategy       string
	Root           string
	NoDeps         bool
	Prefix         string
	IsFileStrategy bool
	File           *path.FilePath
	Output         *path.FilePath
	Symbols        *config.RoutesSymbols
}

type rconfigrule struct {
	fileExts    []string
	outputExts  []string
	strategies  map[string]struct{}
	paramSym    map[string]struct{}
	wildcardSym map[string]struct{}
	rootSymb    map[string]struct{}
}

type rconfig struct {
	cfg  *config.RoutesConfig
	rule *rconfigrule
}

func NewConfig(cfg *config.RoutesConfig) RoutesConfig {
	return &rconfig{
		cfg:  cfg,
		rule: newConfigRule(),
	}
}

func (c *rconfig) Load() (*Config, error) {
	if c.cfg.Strategy == "" {
		msg := fmt.Sprintf("strategy: <empty>, %q must be provided", "strategy")

		return nil, appErr.NewValidationError("strategy", msg)
	}

	if !c.IsValidStartegy() {
		msg := fmt.Sprintf("strategy: %q invalid, allowed: %s",
			c.cfg.Strategy, "file, next-app or next-page")

		return nil, appErr.NewValidationError("strategy", msg)
	}

	if c.IsFileStrategy() && c.cfg.File == "" {
		msg := fmt.Sprintf("strategy: %q, %q must be provided", c.cfg.Strategy, "file")

		return nil, appErr.NewValidationError("file", msg)
	}

	if !c.IsFileStrategy() && c.cfg.File != "" {
		msg := fmt.Sprintf("strategy: %q, %q must be unset", c.cfg.Strategy, "file")

		return nil, appErr.NewValidationError("file", msg)
	}

	output, err := path.ParseFilePath(c.cfg.Output, c.rule.outputExts)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Strategy:       c.cfg.Strategy,
		Root:           c.normalizeRoot(),
		NoDeps:         c.cfg.NoDeps,
		Prefix:         c.cfg.Prefix,
		IsFileStrategy: c.IsFileStrategy(),
		Output:         output,
		Symbols:        c.normalizeSymbols(),
	}

	if c.cfg.File != "" {
		file, err := path.ParseFilePath(c.cfg.File, c.rule.fileExts)
		if err != nil {
			return nil, err
		}
		cfg.File = file
	}

	return cfg, nil
}

func (c *rconfig) IsFileStrategy() bool {
	return c.cfg.Strategy == "file"
}

func (c *rconfig) IsValidStartegy() bool {
	_, ok := c.rule.strategies[c.cfg.Strategy]
	return ok
}

func (c *rconfig) normalizeSymbols() *config.RoutesSymbols {
	if _, ok := c.rule.paramSym[c.cfg.Symbols.Param]; !ok {
		c.cfg.Symbols.Param = "$param"
	}

	if _, ok := c.rule.wildcardSym[c.cfg.Symbols.Wildcard]; !ok {
		c.cfg.Symbols.Wildcard = "$wildcard"
	}

	if _, ok := c.rule.rootSymb[c.cfg.Symbols.Root]; !ok {
		c.cfg.Symbols.Root = "root"
	}

	return c.cfg.Symbols
}

func (c *rconfig) normalizeRoot() string {
	c.cfg.Root = filepath.Clean(c.cfg.Root)

	var suffix string
	switch c.cfg.Strategy {
	case "next-app":
		suffix = "app"

	case "next-page":
		suffix = "pages"

	default:
		suffix = ""
	}

	if suffix == "" {
		return c.cfg.Root
	}

	if filepath.Base(c.cfg.Root) == suffix {
		return c.cfg.Root
	}

	c.cfg.Root = filepath.Join(c.cfg.Root, suffix)

	return c.cfg.Root
}

func newConfigRule() *rconfigrule {
	return &rconfigrule{
		fileExts:   []string{".json", ".yaml", ".yml"},
		outputExts: []string{".js", ".ts"},
		strategies: map[string]struct{}{
			"file":      {},
			"next-app":  {},
			"next-page": {},
		},
		paramSym: map[string]struct{}{
			"$p":     {},
			"$param": {},
		},
		wildcardSym: map[string]struct{}{
			"$w":        {},
			"$wildcard": {},
		},
		rootSymb: map[string]struct{}{
			"p":    {},
			"root": {},
			"path": {},
		},
	}
}
