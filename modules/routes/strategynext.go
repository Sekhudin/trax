package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	appErr "trax/internal/errors"
)

type nextapp struct{}

type nextpage struct{}

type nextrule struct {
	Params affix
	Slug   affix
	OSlug  affix

	Group affix

	DotAffix       affix
	DotDotAffix    affix
	DoubleDotAffix affix
	TripleDotAffix affix

	page walkrule
	app  walkrule
}

var (
	nextExts = map[string]struct{}{
		".js":  {},
		".jsx": {},
		".tsx": {},
	}

	nextRule = nextrule{
		Params: affix{pre: "[", suf: "]"},
		Slug:   affix{pre: "[..", suf: "]"},
		OSlug:  affix{pre: "[[...", suf: "]]"},

		Group: affix{pre: "(", suf: ")"},

		DotAffix:       affix{pre: "(.)"},
		DotDotAffix:    affix{pre: "(..)"},
		DoubleDotAffix: affix{pre: "(..)(..)"},
		TripleDotAffix: affix{pre: "(...)"},

		app: walkrule{
			exts: nextExts,
			identRoute: map[string]struct{}{
				"page": {},
			},

			excludeFiles: map[string]struct{}{
				"api":   {},
				"route": {},
			},
		},

		page: walkrule{
			exts: nextExts,
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
		},
	}
)

func (*nextapp) shouldSkip(p string, d fs.DirEntry) error {
	return nextRule.shouldSkip(d)
}

func (*nextapp) normalizeSegment(seg string) (string, error) {
	return nextRule.normalizeSegment(seg)
}

func (*nextpage) shouldSkip(p string, d fs.DirEntry) error {
	return nextRule.shouldSkip(d)
}

func (*nextpage) normalizeSegment(seg string) (string, error) {
	return nextRule.normalizeSegment(seg)
}

func (n *nextrule) normalizeSegment(seg string) (string, error) {
	if seg == "" {
		return seg, nil
	}

	kind := n.segmentKind(seg)
	if kind == "Static" {
		return seg, nil
	}

	a, err := n.getAffix(kind)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(seg, a.pre) {
		if !strings.HasSuffix(seg, a.suf) {
			msg := fmt.Sprintf("%q invalid segment", seg)

			return "", appErr.NewInvalidConfigError("path", msg)
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

func (n *nextrule) segmentKind(seg string) string {
	switch {
	case strings.HasPrefix(seg, n.OSlug.pre):
		return "OSlug"

	case strings.HasPrefix(seg, n.Slug.pre):
		return "Slug"

	case strings.HasPrefix(seg, n.Params.pre):
		return "Params"

	default:
		return "Static"
	}
}

func (n *nextrule) getAffix(field string) (affix, error) {
	rv := reflect.ValueOf(n)
	rv = reflect.Indirect(rv).FieldByName(field)

	if !rv.IsValid() {
		return affix{}, appErr.NewInvalidConfigError("path", "affix struct not found")
	}

	return rv.Interface().(affix), nil
}

func (*nextrule) removeAffix(seg string, a affix) string {
	seg = strings.TrimPrefix(seg, a.pre)
	seg = strings.TrimSuffix(seg, a.suf)

	return seg
}

func (n *nextrule) shouldSkip(d fs.DirEntry) error {
	if d.IsDir() {
		dir := d.Name()
		if strings.HasPrefix(dir, n.DotAffix.pre) ||
			strings.HasPrefix(dir, n.DotDotAffix.pre) ||
			strings.HasPrefix(dir, n.DoubleDotAffix.pre) ||
			strings.HasPrefix(dir, n.TripleDotAffix.pre) {

			return filepath.SkipDir
		}
	}

	return nil
}
