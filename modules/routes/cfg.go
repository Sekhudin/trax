package routes

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"trax/internal/path"

	"github.com/spf13/viper"
)

type Cfg struct {
	Strategy string
	Root     string
	File     *path.FilePath
	Output   *path.FilePath
	Oext     string
}

var (
	rExts      = []string{".json", ".yaml", ".yml"}
	oExts      = []string{".js", ".ts"}
	strategies = []string{"file", "next-app", "next-page"}
)

func NewCfg() (*Cfg, error) {
	cfg := Cfg{
		Strategy: viper.GetString("routes.strategy"),
		Root:     viper.GetString("routes.root"),
	}

	file := viper.GetString("routes.file")
	output := viper.GetString("routes.output")

	if cfg.Strategy == "" {
		return nil, fmt.Errorf("strategy: <empty>, %q must be provided", "strategy")
	}

	if cfg.Strategy == "file" && file == "" {
		return nil, fmt.Errorf("strategy: %q, %q must be provided", cfg.Strategy, "file")
	}

	if cfg.Strategy != "file" && file != "" {
		return nil, fmt.Errorf("strategy: %q, %q must be unset", cfg.Strategy, "file")
	}

	if !cfg.isValidStartegy() {
		return nil, fmt.Errorf("strategy: %q invalid, allowed: %q",
			cfg.Strategy,
			strings.Join(strategies, " | "))
	}

	if file != "" {
		oPath, err := path.ParseFilePath(file, rExts)
		if err != nil {
			return nil, err
		}
		cfg.File = oPath
	}

	oPath, err := path.ParseFilePath(output, oExts)
	if err != nil {
		return nil, err
	}

	cfg.Root = cfg.normalizeRoot()
	cfg.Output = oPath

	return &cfg, nil
}

func (c *Cfg) isValidStartegy() bool {
	return slices.Contains(strategies, c.Strategy)
}

func (c *Cfg) normalizeRoot() string {
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
