package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
)

type stgNext struct{}

type stgNextPage struct{}

type stgNextApp struct{}

type stgNextAffixRule struct {
	Params affix
	Slug   affix
	OSlug  affix

	Group affix

	DotAffix       affix
	DotDotAffix    affix
	DoubleDotAffix affix
	TripleDotAffix affix
}

var stgNextAffix = stgNextAffixRule{
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
	return stgNextShouldSkip(d)
}

func (stgNextPage) normalizeSegment(seg string) (string, error) {
	return stgNextNormalizeSegment(seg)
}

func (stgNextApp) shouldSkip(p string, d fs.DirEntry) error {
	return stgNextShouldSkip(d)
}

func (stgNextApp) normalizeSegment(seg string) (string, error) {
	return stgNextNormalizeSegment(seg)
}

func stgNextNormalizeSegment(seg string) (string, error) {
	if seg == "" {
		return seg, nil
	}

	kind := stgNextSegmentKind(seg)
	if kind == "Static" {
		return seg, nil
	}

	a, err := stgNextGetAffix(kind)
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

func stgNextSegmentKind(seg string) string {
	switch {
	case strings.HasPrefix(seg, stgNextAffix.OSlug.pre):
		return "OSlug"

	case strings.HasPrefix(seg, stgNextAffix.Slug.pre):
		return "Slug"

	case strings.HasPrefix(seg, stgNextAffix.Params.pre):
		return "Params"

	default:
		return "Static"
	}
}

func stgNextGetAffix(field string) (affix, error) {
	r := reflect.ValueOf(stgNextAffix)
	f := reflect.Indirect(r).FieldByName(field)

	if !f.IsValid() {
		return affix{}, fmt.Errorf("affix struct not found")
	}

	return f.Interface().(affix), nil
}

func stgNextRemoveAffix(seg string, a affix) string {
	seg = strings.TrimPrefix(seg, a.pre)
	seg = strings.TrimSuffix(seg, a.suf)

	return seg
}

func stgNextShouldSkip(d fs.DirEntry) error {
	if d.IsDir() {
		dir := d.Name()
		if strings.HasPrefix(dir, stgNextAffix.DotAffix.pre) ||
			strings.HasPrefix(dir, stgNextAffix.DotDotAffix.pre) ||
			strings.HasPrefix(dir, stgNextAffix.DoubleDotAffix.pre) ||
			strings.HasPrefix(dir, stgNextAffix.TripleDotAffix.pre) {

			return filepath.SkipDir
		}
	}

	return nil
}
