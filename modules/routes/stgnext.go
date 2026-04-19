package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
)

type stgNextPage struct{}

type stgNextApp struct{}

type stgNextRule struct {
	Params affix
	Slug   affix
	OSlug  affix

	Group affix

	DotAffix       affix
	DotDotAffix    affix
	DoubleDotAffix affix
	TripleDotAffix affix
}

var nextRule = stgNextRule{
	Params: affix{pre: "[", suf: "]"},
	Slug:   affix{pre: "[..", suf: "]"},
	OSlug:  affix{pre: "[[...", suf: "]]"},

	Group: affix{pre: "(", suf: ")"},

	DotAffix:       affix{pre: "(.)"},
	DotDotAffix:    affix{pre: "(..)"},
	DoubleDotAffix: affix{pre: "(..)(..)"},
	TripleDotAffix: affix{pre: "(...)"},
}

var nextPageRule = rule{
	exts: map[string]struct{}{
		".js":  {},
		".jsx": {},
		".tsx": {},
	},

	identRoute: map[string]struct{}{
		"index": {},
	},

	excludeFiles: map[string]struct{}{
		"_app":      {},
		"_document": {},
		"_error":    {},
		"api":       {},
		"404":       {},
		"500":       {},
	},
}

var nextAppRule = rule{
	exts: nextPageRule.exts,

	identRoute: map[string]struct{}{
		"page": {},
	},

	excludeFiles: map[string]struct{}{
		"api":   {},
		"route": {},
	},
}

func (stgNextPage) shouldSkip(p string, d fs.DirEntry) error {
	return nextRule.shouldSkip(d)
}

func (stgNextPage) normalizeSegment(seg string) (string, error) {
	return nextRule.normalizeSegment(seg)
}

func (stgNextApp) shouldSkip(p string, d fs.DirEntry) error {
	return nextRule.shouldSkip(d)
}

func (stgNextApp) normalizeSegment(seg string) (string, error) {
	return nextRule.normalizeSegment(seg)
}

func (r *stgNextRule) normalizeSegment(seg string) (string, error) {
	if seg == "" {
		return seg, nil
	}

	kind := r.segmentKind(seg)
	if kind == "Static" {
		return seg, nil
	}

	a, err := r.getAffix(kind)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(seg, a.pre) {
		if !strings.HasSuffix(seg, a.suf) {
			return "", fmt.Errorf("%q invalid segment", seg)
		}
	}

	seg = strings.TrimPrefix(seg, a.pre)
	seg = strings.TrimSuffix(seg, a.suf)
	seg = strings.TrimSpace(seg)

	switch kind {
	case "Params":
		return fmt.Sprintf(":%s", seg), nil

	case "Slug", "OSlug":
		return "*", nil

	default:
		return seg, nil
	}
}

func (r *stgNextRule) segmentKind(seg string) string {
	switch {
	case strings.HasPrefix(seg, r.OSlug.pre):
		return "OSlug"

	case strings.HasPrefix(seg, r.Slug.pre):
		return "Slug"

	case strings.HasPrefix(seg, r.Params.pre):
		return "Params"

	default:
		return "Static"
	}
}

func (r *stgNextRule) getAffix(field string) (affix, error) {
	rv := reflect.ValueOf(r)
	rv = reflect.Indirect(rv).FieldByName(field)

	if !rv.IsValid() {
		return affix{}, fmt.Errorf("affix struct not found")
	}

	return rv.Interface().(affix), nil
}

func (r *stgNextRule) removeAffix(seg string, a affix) string {
	seg = strings.TrimPrefix(seg, a.pre)
	seg = strings.TrimSuffix(seg, a.suf)

	return seg
}

func (r *stgNextRule) shouldSkip(d fs.DirEntry) error {
	if d.IsDir() {
		dir := d.Name()
		if strings.HasPrefix(dir, r.DotAffix.pre) ||
			strings.HasPrefix(dir, r.DotDotAffix.pre) ||
			strings.HasPrefix(dir, r.DoubleDotAffix.pre) ||
			strings.HasPrefix(dir, r.TripleDotAffix.pre) {

			return filepath.SkipDir
		}
	}

	return nil
}
