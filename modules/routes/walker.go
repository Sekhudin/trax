package routes

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type affix struct {
	pre string
	suf string
}

type walkrule struct {
	exts         map[string]struct{}
	identRoute   map[string]struct{}
	excludeFiles map[string]struct{}
}

type walkerstrategy interface {
	shouldSkip(p string, d fs.DirEntry) error
	normalizeSegment(seg string) (string, error)
}

type walker struct {
	strategy walkerstrategy
	config   *Config
	rule     *walkrule
}

func (w *walker) walk() ([]RawRoute, error) {
	var rs []RawRoute

	c := *w.config
	ws := w.strategy

	err := filepath.WalkDir(c.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return appErr.NewIOError("path", fmt.Sprintf("failed to access path %q", p), err)
		}

		if err := ws.shouldSkip(p, d); err != nil {
			return err
		}

		if !w.isRouteFile(p, d) {
			return nil
		}

		rel, err := filepath.Rel(c.Root, p)
		if err != nil {
			return appErr.NewInternalError("path", fmt.Sprintf("failed to calculate relative path %q", p), err)
		}

		segs := w.splitSegments(rel)
		name := w.buildName(segs)

		segs, err = w.normalizeSegments(segs)
		if err != nil {
			return nil
		}

		path := w.buildPath(segs)

		rs = append(rs, RawRoute{Name: name, Path: path})

		return nil
	})

	return rs, err
}

func (w *walker) isRouteFile(p string, d fs.DirEntry) bool {
	if d.IsDir() {
		return false
	}

	r := w.rule

	ext := filepath.Ext(p)

	if _, ok := r.exts[ext]; !ok {
		return false
	}

	base := strings.TrimSuffix(d.Name(), ext)
	if _, ok := r.excludeFiles[base]; ok {
		return false
	}

	return true
}

func (w *walker) splitSegments(rel string) []string {
	res := []string{}

	r := w.rule

	ext := filepath.Ext(rel)
	base := strings.TrimSuffix(rel, ext)

	if _, ok := r.identRoute[base]; ok {
		res := append(res, "$root")

		return res
	}

	segs := strings.SplitSeq(base, string(filepath.Separator))
	for seg := range segs {
		if _, ok := r.identRoute[seg]; ok {
			continue
		}

		res = append(res, seg)
	}

	return res
}

func (w *walker) normalizeSegments(segs []string) ([]string, error) {
	res := []string{}

	ws := w.strategy

	for _, seg := range segs {
		normal, err := ws.normalizeSegment(seg)
		if err != nil {
			return nil, err
		}

		res = append(res, normal)
	}

	return res, nil
}

func (*walker) buildPath(segs []string) string {
	for i, seg := range segs {
		if seg == "$root" {
			segs[i] = ""
			continue
		}
	}

	p := fmt.Sprintf("/%s", strings.Join(segs, "/"))
	p = path.Clean(p)

	return p
}

func (w *walker) buildName(parts []string) string {
	var words []string

	for _, p := range parts {
		clean := w.sanitizeInput(p)
		if clean == "" {
			continue
		}

		subs := strings.FieldsFunc(clean, func(r rune) bool {
			return r == '-' || r == '_'
		})
		words = append(words, subs...)
	}

	return w.toSnakeCase(words)
}

func (w *walker) sanitizeInput(s string) string {
	nonAlnum := regexp.MustCompile(`[^a-zA-Z0-9-]+`)

	s = strings.TrimSpace(s)
	return nonAlnum.ReplaceAllString(s, "")
}

func (w *walker) toSnakeCase(words []string) string {
	if len(words) == 0 {
		return ""
	}

	for i, w := range words {
		words[i] = strings.ToLower(w)
	}

	return strings.Join(words, "_")
}
