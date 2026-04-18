package routes

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"trax/internal/path"

	"github.com/spf13/viper"
)

type RoutesConfig struct {
	Strategy string
	Root     string
	File     *path.FilePath
	Output   *path.FilePath
}

var (
	prefRoute = "routes"
	identRgx  = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	staticRgx = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

	fileExts   = []string{".json", ".yaml", ".yml"}
	outputExts = []string{".js", ".ts"}
	strategies = []string{"file", "next-app", "next-page"}
)

func LoadConfig() (*RoutesConfig, error) {
	cfg := RoutesConfig{
		Strategy: viper.GetString("routes.strategy"),
		Root:     viper.GetString("routes.root"),
	}

	file := viper.GetString("routes.file")
	output := viper.GetString("routes.output")

	if cfg.Strategy == "" {
		return nil, fmt.Errorf("strategy: <empty>, 'strategy' must be provided")
	}

	if cfg.Strategy == "file" && file == "" {
		return nil, fmt.Errorf("strategy: '%s', 'file' must be provided", cfg.Strategy)
	}

	if cfg.Strategy != "file" && file != "" {
		return nil, fmt.Errorf("strategy: '%s', 'file' must be unset", cfg.Strategy)
	}

	if !isAllowed(cfg.Strategy, strategies) {
		return nil, fmt.Errorf("strategy: '%s' invalid, allowed: %s",
			cfg.Strategy,
			strings.Join(strategies, ", "))
	}

	if file != "" {
		oPath, err := path.ParseFilePath(file, fileExts)
		if err != nil {
			return nil, err
		}
		cfg.File = oPath
	}

	oPath, err := path.ParseFilePath(output, outputExts)
	if err != nil {
		return nil, err
	}

	cfg.Root = filepath.Clean(cfg.Root)
	cfg.Output = oPath

	return &cfg, nil
}

func isAllowed(ext string, allowed []string) bool {
	return slices.Contains(allowed, ext)
}
