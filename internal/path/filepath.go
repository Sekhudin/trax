package path

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

type FilePath struct {
	Dir      string
	Filename string
	Ext      string
	Full     string
}

func ParseFilePath(out string, allowedExts []string) (*FilePath, error) {
	clean := strings.TrimSpace(out)
	if clean == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	ext := filepath.Ext(clean)
	if ext == "" {
		return nil, fmt.Errorf(
			"path: '%s' must include file extension (allowed: %s)",
			clean,
			strings.Join(allowedExts, ", "),
		)
	}

	if !isAllowedExt(ext, allowedExts) {
		return nil, fmt.Errorf(
			"unsupported extension: %s (allowed: %s)",
			ext,
			strings.Join(allowedExts, ", "),
		)
	}

	return &FilePath{
		Full:     clean,
		Dir:      filepath.Dir(clean),
		Filename: filepath.Base(clean),
		Ext:      ext,
	}, nil
}

func isAllowedExt(ext string, allowed []string) bool {
	return slices.Contains(allowed, ext)
}
