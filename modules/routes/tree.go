package routes

import (
	"fmt"
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type treeselector func(selector string) (map[string]any, error)

type treebuilderItf interface {
	build([]route) (map[string]*node, treeselector, error)
}

type treebuilder struct {
	cfg *Config
}

func newTreeBuilder(cfg *Config) treebuilderItf {
	return &treebuilder{cfg: cfg}
}

func (b *treebuilder) build(rs []route) (map[string]*node, treeselector, error) {
	tr := make(map[string]*node)

	for _, r := range rs {
		if err := b.insert(r, tr); err != nil {
			return nil, nil, err
		}
	}

	trs := b.createSelector(b.toMap(tr))

	return tr, trs, nil
}

func (*treebuilder) createSelector(tr map[string]any) treeselector {
	prefix := viper.GetString("routes.prefix")

	v := viper.New()

	v.SetDefault(prefix, tr[prefix])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", prefix, selector)
		} else {
			selector = prefix
		}

		selector = strings.TrimSuffix(selector, "$")
		selector = strings.TrimSuffix(selector, ".")
		selector = strings.TrimSuffix(selector, "?")
		selector = strings.ReplaceAll(selector, "?.", ".$")
		val := v.Get(selector)

		switch v := val.(type) {

		case nil:
			msg := fmt.Sprintf("%q not found", selector)

			return nil, appErr.NewValidationError("selector", msg)

		case map[string]any:
			return v, nil

		case string:
			return map[string]any{
				"path": v,
			}, nil

		default:
			msg := fmt.Sprintf("invalid type for selector %q: %T",
				selector, val,
			)

			return nil, appErr.NewValidationError("selector", msg)
		}
	}
}

func (b *treebuilder) toMap(nds map[string]*node) map[string]any {
	keys := make([]string, 0, len(nds))
	for k := range nds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make(map[string]any, len(nds))

	for _, k := range keys {
		n := nds[k]

		if n == nil {
			continue
		}

		m := make(map[string]any)

		if n.root != "" {
			m[b.cfg.Symbols.Root] = n.root
		}

		if len(n.children) > 0 {
			maps.Copy(m, b.toMap(n.children))
		}

		result[k] = m
	}

	return result
}

func (b *treebuilder) insert(r route, current map[string]*node) error {
	for i, part := range r.parts {
		key, kind, err := b.normalizePart(r, part)
		if err != nil {
			return err
		}

		nd, ok := current[key]
		if !ok {
			if err := b.validateChild(r, current, kind); err != nil {
				return err
			}
			nd = &node{
				children: make(map[string]*node),
				kind:     kind,
				segment:  part,
			}
			current[key] = nd
		} else if nd.segment != part {
			msg := fmt.Sprintf("conflicting segment %q and %q produce same key", nd.segment, part)
			return appErr.NewValidationError("path", msg)
		}

		if i == len(r.parts)-1 {
			if nd.root != "" && nd.root != r.path {
				msg := fmt.Sprintf("%q duplicate route detected: %q vs %q", r.name, nd.root, r.path)
				return appErr.NewValidationError("path", msg)
			}
			nd.root = r.path
		}

		current = nd.children
	}

	return nil
}

func (b *treebuilder) normalizePart(r route, part string) (string, string, error) {
	identPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	staticPattern := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

	if part == "*" {
		return b.cfg.Symbols.Wildcard, "wildcard", nil
	}

	if cut, found := strings.CutPrefix(part, ":"); found {
		if !identPattern.MatchString(cut) {
			msg := fmt.Sprintf("%q invalid param name: %s", r.name, cut)

			return "", "", appErr.NewValidationError("path", msg)
		}
		return b.cfg.Symbols.Param, "param", nil
	}

	if !staticPattern.MatchString(part) {
		msg := fmt.Sprintf("%q invalid path segment: %s", r.name, part)

		return "", "", appErr.NewValidationError("path", msg)
	}

	if strings.Contains(part, "-") {
		parts := strings.Split(part, "-")
		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, ""), "static", nil
	}

	return part, "static", nil
}

func (b *treebuilder) validateChild(r route, current map[string]*node, kind string) error {
	for _, c := range current {
		if kind == "wildcard" || c.kind == "wildcard" {
			msg := fmt.Sprintf("%q path wildcard route cannot coexist with other routes at the same level", r.name)

			return appErr.NewValidationError("path", msg)
		}

		if kind == "param" && c.kind == "static" {
			msg := fmt.Sprintf("%q path param route conflicts with existing static route", r.name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "static" && c.kind == "param" {
			msg := fmt.Sprintf("%q path static route conflicts with existing param route", r.name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "param" && c.kind == "param" {
			msg := fmt.Sprintf("%q path multiple param routes at the same level are not allowed", r.name)

			return appErr.NewValidationError("path", msg)
		}
	}

	return nil
}
