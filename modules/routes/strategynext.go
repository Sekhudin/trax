package routes

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type nextrule struct {
	NonRouteDir affix
	SlotDir     affix

	Group  affix
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

type nextapp struct {
	rule *nextrule
}

type nextpage struct {
	rule *nextrule
}

func newNextApp() *nextapp {
	return &nextapp{
		rule: newNextRule(),
	}
}

func newNextPage() *nextpage {
	return &nextpage{
		rule: newNextRule(),
	}
}

func (n *nextapp) shouldSkip(p string, d fs.DirEntry) error {
	if n.rule.isNonRouteDir(d) {
		return filepath.SkipDir
	}

	if n.rule.isSlotDir(d) {
		return filepath.SkipDir
	}

	return n.rule.shouldSkip(d)
}

func (n *nextapp) normalizeSegment(seg string) string {
	return n.rule.normalizeSegment(seg)
}

func (n *nextpage) shouldSkip(p string, d fs.DirEntry) error {
	return n.rule.shouldSkip(d)
}

func (n *nextpage) normalizeSegment(seg string) string {
	return n.rule.normalizeSegment(seg)
}

func newNextRule() *nextrule {
	nextExts := map[string]struct{}{
		".js":  {},
		".jsx": {},
		".tsx": {},
	}

	return &nextrule{
		NonRouteDir: affix{pre: "_"},
		SlotDir:     affix{pre: "@"},

		Group:  affix{pre: "(", suf: ")"},
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
}

func (n *nextrule) isNonRouteDir(d fs.DirEntry) bool {
	if d.IsDir() {
		return strings.HasPrefix(d.Name(), n.NonRouteDir.pre)
	}

	return false
}

func (n *nextrule) isSlotDir(d fs.DirEntry) bool {
	if d.IsDir() {
		return strings.HasPrefix(d.Name(), n.SlotDir.pre)
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

func (n *nextrule) normalizeSegment(seg string) string {
	if seg == "" {
		return seg
	}

	kind := n.segmentKind(seg)
	if kind == "Static" {
		r := strings.NewReplacer(
			"(", "p_",
			")", "_s",
		)

		seg = r.Replace(seg)
	}

	r := strings.NewReplacer(
		".", "",
		"@", "",
		" ", "",
		"[", "",
		"]", "",
		"(", "",
		")", "",
	)

	seg = r.Replace(seg)

	switch kind {
	case "Group":
		return ""

	case "Params":
		return fmt.Sprintf(":%s", seg)

	case "Slug", "OSlug":
		return "*"

	default:
		return seg
	}
}

func (n *nextrule) segmentKind(seg string) string {
	switch {
	case strings.HasPrefix(seg, n.Group.pre) && strings.HasSuffix(seg, n.Group.suf):
		return "Group"

	case strings.HasPrefix(seg, n.OSlug.pre) && strings.HasSuffix(seg, n.OSlug.suf):
		return "OSlug"

	case strings.HasPrefix(seg, n.Slug.pre) && strings.HasSuffix(seg, n.Slug.suf):
		return "Slug"

	case strings.HasPrefix(seg, n.Params.pre) && strings.HasSuffix(seg, n.Params.suf):
		return "Params"

	default:
		return "Static"
	}
}
