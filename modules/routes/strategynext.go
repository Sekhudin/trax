package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type nextapp struct{}

type nextpage struct{}

type nextrule struct {
	NonRouteDir affix
	SlotDir     affix
	GroupDir    affix

	Params affix
	Slug   affix
	OSlug  affix

	DotAffix       affix
	DotDotAffix    affix
	DoubleDotAffix affix
	TripleDotAffix affix

	skipFolders map[string]struct{}

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
		NonRouteDir: affix{pre: "_"},
		SlotDir:     affix{pre: "@"},
		GroupDir:    affix{pre: "(", suf: ")"},

		Params: affix{pre: "[", suf: "]"},
		Slug:   affix{pre: "[...", suf: "]"},
		OSlug:  affix{pre: "[[...", suf: "]]"},

		DotAffix:       affix{pre: "(.)"},
		DotDotAffix:    affix{pre: "(..)"},
		DoubleDotAffix: affix{pre: "(..)(..)"},
		TripleDotAffix: affix{pre: "(...)"},

		skipFolders: map[string]struct{}{
			"api": {},
		},

		app: walkrule{
			exts: nextExts,
			identRoute: map[string]struct{}{
				"page": {},
			},

			excludeFiles: map[string]struct{}{
				"instrumentation": {},
				"proxy":           {},
				"api":             {},

				"layout":       {},
				"loading":      {},
				"not-found":    {},
				"error":        {},
				"global-error": {},
				"route":        {},
				"template":     {},
				"default":      {},

				"favicon":    {},
				"icon":       {},
				"apple-icon": {},

				"opengraph-image": {},
				"twitter-image":   {},

				"sitemap": {},
				"robots":  {},
			},
		},

		page: walkrule{
			exts: nextExts,
			identRoute: map[string]struct{}{
				"index": {},
			},

			excludeFiles: map[string]struct{}{
				"instrumentation": {},
				"proxy":           {},
				"api":             {},

				"_app":      {},
				"_document": {},
				"_error":    {},
				"404":       {},
				"500":       {},
			},
		},
	}
)

func (*nextapp) shouldSkip(p string, d fs.DirEntry) error {
	if nextRule.isNonRouteDir(d) {
		return filepath.SkipDir
	}

	if nextRule.isSlotDir(d) {
		return filepath.SkipDir
	}

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

func (n *nextrule) isNonRouteDir(d fs.DirEntry) bool {
	if d.IsDir() {
		return strings.HasPrefix(d.Name(), nextRule.NonRouteDir.pre)
	}

	return false
}

func (n *nextrule) isSlotDir(d fs.DirEntry) bool {
	if d.IsDir() {
		return strings.HasPrefix(d.Name(), nextRule.SlotDir.pre)
	}

	return false
}

func (n *nextrule) shouldSkip(d fs.DirEntry) error {
	if d.IsDir() {
		dir := d.Name()

		if _, ok := n.skipFolders[dir]; ok {
			return filepath.SkipDir
		}

		if strings.HasPrefix(dir, n.DotAffix.pre) ||
			strings.HasPrefix(dir, n.DotDotAffix.pre) ||
			strings.HasPrefix(dir, n.DoubleDotAffix.pre) ||
			strings.HasPrefix(dir, n.TripleDotAffix.pre) {

			return filepath.SkipDir
		}
	}

	return nil
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
	case "Group":
		return "", nil

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
	case strings.HasPrefix(seg, n.GroupDir.pre) && strings.HasSuffix(seg, n.GroupDir.suf):
		return "Group"

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
