package routes

import (
	"fmt"
	"path/filepath"

	"github.com/sekhudin/trax/internal/path"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type symbol struct {
	Param    string
	Wildcard string
	Root     string
}

type Config struct {
	Strategy       string
	IsFileStrategy bool
	Root           string
	NoDeps         bool
	File           *path.FilePath
	Output         *path.FilePath
	Oext           string
	Symbols        *symbol
}

type configrule struct {
	fileExts    []string
	outputExts  []string
	strategies  map[string]struct{}
	paramSym    map[string]struct{}
	wildcardSym map[string]struct{}
	rootSymb    map[string]struct{}
}

var cfgRule = newConfigRule()

func NewConfig() (*Config, error) {
	cfg := Config{
		Strategy: viper.GetString("routes.strategy"),
		Root:     viper.GetString("routes.root"),
		NoDeps:   viper.GetBool("routes.no-deps"),
	}

	cfg.IsFileStrategy = cfgRule.IsFileStrategy(cfg.Strategy)

	file := viper.GetString("routes.file")
	output := viper.GetString("routes.output")

	if cfg.Strategy == "" {
		msg := fmt.Sprintf("strategy: <empty>, %q must be provided", "strategy")

		return nil, appErr.NewValidationError("strategy", msg)
	}

	if !cfgRule.isValidStartegy(cfg.Strategy) {
		msg := fmt.Sprintf("strategy: %q invalid, allowed: %s",
			cfg.Strategy, "file, next-app or next-page")

		return nil, appErr.NewValidationError("strategy", msg)
	}

	if cfg.IsFileStrategy && file == "" {
		msg := fmt.Sprintf("strategy: %q, %q must be provided", cfg.Strategy, "file")

		return nil, appErr.NewValidationError("file", msg)
	}

	if !cfg.IsFileStrategy && file != "" {
		msg := fmt.Sprintf("strategy: %q, %q must be unset", cfg.Strategy, "file")

		return nil, appErr.NewValidationError("file", msg)
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

	cfg.Symbols = cfgRule.normalizeSymbols()
	cfg.Root = cfgRule.normalizeRootPath(cfg.Root, cfg.Strategy)
	cfg.Output = oPath

	return &cfg, nil
}

func newConfigRule() *configrule {
	return &configrule{
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
			"$":    {},
			"r":    {},
			"p":    {},
			"root": {},
			"path": {},
		},
	}
}

func (*configrule) IsFileStrategy(strategy string) bool {
	return strategy == "file"
}

func (*configrule) isValidStartegy(strategy string) bool {
	_, ok := cfgRule.strategies[strategy]
	return ok
}

func (r *configrule) normalizeSymbols() *symbol {
	sym := symbol{
		Param:    viper.GetString("routes.symbols.param"),
		Wildcard: viper.GetString("routes.symbols.wildcard"),
		Root:     viper.GetString("routes.symbols.root"),
	}

	if _, ok := r.paramSym[sym.Param]; !ok {
		sym.Param = "$param"
	}

	if _, ok := r.paramSym[sym.Wildcard]; !ok {
		sym.Wildcard = "$wildcard"
	}

	if _, ok := r.paramSym[sym.Root]; !ok {
		sym.Root = "$root"
	}

	return &sym
}

func (r *configrule) normalizeRootPath(root string, strategy string) string {
	root = filepath.Clean(root)

	var suffix string
	switch strategy {
	case "next-app":
		suffix = "app"

	case "next-page":
		suffix = "pages"

	default:
		suffix = ""
	}

	if suffix == "" {
		return root
	}

	if filepath.Base(root) == suffix {
		return root
	}

	root = filepath.Join(root, suffix)
	return root
}
