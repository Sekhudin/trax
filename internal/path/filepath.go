package path

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type FilePath struct {
	Dir      string
	Filename string
	Ext      string
	Full     string
}

var filepathRel = filepath.Rel

func ParseFilePath(out string, allowedExts []string) (*FilePath, error) {
	clean := strings.TrimSpace(out)
	if clean == "" {
		return nil, appErr.NewValidationError("path", "path cannot be empty")
	}

	rel, err := filepathRel(".", clean)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(rel))
	if ext == "" {
		msg := fmt.Sprintf(
			"path: %q must include file extension (allowed: %q)",
			clean, strings.Join(allowedExts, " | "),
		)
		return nil, appErr.NewValidationError("path_format", msg)
	}

	if !isAllowedExt(ext, allowedExts) {
		msg := fmt.Sprintf(
			"unsupported extension: %s (allowed: %q)",
			ext,
			strings.Join(allowedExts, " | "),
		)
		return nil, appErr.NewValidationError("extension", msg)
	}

	return &FilePath{
		Full:     rel,
		Dir:      filepath.Dir(rel),
		Filename: filepath.Base(rel),
		Ext:      ext,
	}, nil
}

func isAllowedExt(ext string, allowed []string) bool {
	return slices.Contains(allowed, ext)
}
