package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type affix struct {
	pre string
	suf string
}

type rule struct {
	exts         map[string]struct{}
	identRoute   map[string]struct{}
	excludeFiles map[string]struct{}
}

type walkRule interface {
	shouldSkip(p string, d fs.DirEntry) error
	normalizeSegment(seg string) (string, error)
}

type walker struct {
	cfg   *Cfg
	rule  *rule
	wRule walkRule
}

func (w *walker) walk() ([]raw, error) {
	var rs []raw

	c := w.cfg
	wr := w.wRule

	err := filepath.WalkDir(c.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err := wr.shouldSkip(p, d); err != nil {
			return err
		}

		if !w.isRouteFile(p, d) {
			return nil
		}

		rel, err := filepath.Rel(c.Root, p)
		if err != nil {
			return err
		}

		segs := w.splitSegments(rel)
		name := w.buildName(segs)

		segs, err = w.normalizeSegments(segs)
		if err != nil {
			return nil
		}

		path := w.buildPath(segs)

		rs = append(rs, raw{Name: name, Path: path})

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
		res := append(res, "root")

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

	wr := w.wRule

	for _, seg := range segs {
		normal, err := wr.normalizeSegment(seg)
		if err != nil {
			return nil, err
		}

		res = append(res, normal)
	}

	return res, nil
}

func (*walker) buildPath(segs []string) string {
	for i, seg := range segs {
		if seg == "root" {
			segs[i] = ""
		}
	}

	return fmt.Sprintf("/%s", strings.Join(segs, "/"))
}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

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
