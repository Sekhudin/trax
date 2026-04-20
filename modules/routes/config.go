package routes

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"trax/internal/path"

	"github.com/spf13/viper"
)

type Config struct {
	Strategy string
	Root     string
	File     *path.FilePath
	Output   *path.FilePath
	Oext     string
}

type configrule struct {
	fileExts   []string
	outputExts []string
	strategies []string
}

var cfgRule = configrule{
	fileExts:   []string{".json", ".yaml", ".yml"},
	outputExts: []string{".js", ".ts"},
	strategies: []string{"file", "next-app", "next-page"},
}

func NewConfig() (*Config, error) {
	cfg := Config{
		Strategy: viper.GetString("routes.strategy"),
		Root:     viper.GetString("routes.root"),
	}

	file := viper.GetString("routes.file")
	output := viper.GetString("routes.output")

	if cfg.Strategy == "" {
		return nil, fmt.Errorf("strategy: <empty>, %q must be provided", "strategy")
	}

	if cfg.IsFileStrategy() && file == "" {
		return nil, fmt.Errorf("strategy: %q, %q must be provided", cfg.Strategy, "file")
	}

	if !cfg.IsFileStrategy() && file != "" {
		return nil, fmt.Errorf("strategy: %q, %q must be unset", cfg.Strategy, "file")
	}

	if !cfg.isValidStartegy() {
		return nil, fmt.Errorf("strategy: %q invalid, allowed: %q",
			cfg.Strategy,
			strings.Join(cfgRule.strategies, " | "))
	}

	if file != "" {
		oPath, err := path.ParseFilePath(file, cfgRule.fileExts)
		if err != nil {
			return nil, err
		}
		cfg.File = oPath
	}

	oPath, err := path.ParseFilePath(output, cfgRule.outputExts)
	if err != nil {
		return nil, err
	}

	cfg.Root = cfg.normalizeRoot()
	cfg.Output = oPath

	return &cfg, nil
}

func (c *Config) IsFileStrategy() bool {
	return c.Strategy == "file"
}

func (c *Config) isValidStartegy() bool {
	return slices.Contains(cfgRule.strategies, c.Strategy)
}

func (c *Config) normalizeRoot() string {
	c.Root = filepath.Clean(c.Root)

	var suffix string
	switch c.Strategy {
	case "next-app":
		suffix = "app"

	case "next-page":
		suffix = "pages"

	default:
		return c.Root
	}

	if filepath.Base(c.Root) == suffix {
		return c.Root
	}

	c.Root = filepath.Join(c.Root, suffix)
	return c.Root
}
