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

type TreeSelector func(selector string) (map[string]any, error)

type TreeBuilder interface {
	Build([]Route) (map[string]*Node, TreeSelector, error)
}

type tree struct {
	cfg *Config
}

func NewTreeBuilder(cfg *Config) TreeBuilder {
	return &tree{cfg: cfg}
}

func (b *tree) Build(rs []Route) (map[string]*Node, TreeSelector, error) {
	tr := make(map[string]*Node)

	for _, r := range rs {
		if err := b.insert(r, tr); err != nil {
			return nil, nil, err
		}
	}

	trs := b.createSelector(b.toMap(tr))

	return tr, trs, nil
}

func (b *tree) createSelector(tr map[string]any) TreeSelector {
	v := viper.New()

	v.SetDefault(b.cfg.Prefix, tr[b.cfg.Prefix])

	return func(selector string) (map[string]any, error) {
		if selector != "" {
			selector = fmt.Sprintf("%s.%s", b.cfg.Prefix, selector)
		} else {
			selector = b.cfg.Prefix
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

func (b *tree) toMap(nds map[string]*Node) map[string]any {
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

		if n.Root != "" {
			m[b.cfg.Symbols.Root] = n.Root
		}

		if len(n.Children) > 0 {
			maps.Copy(m, b.toMap(n.Children))
		}

		result[k] = m
	}

	return result
}

func (b *tree) insert(r Route, current map[string]*Node) error {
	for i, part := range r.Parts {
		key, Kind, err := b.normalizePart(r, part)
		if err != nil {
			return err
		}

		nd, ok := current[key]
		if !ok {
			if err := b.validateChild(r, current, Kind); err != nil {
				return err
			}
			nd = &Node{
				Children: make(map[string]*Node),
				Kind:     Kind,
				Segment:  part,
			}
			current[key] = nd
		} else if nd.Segment != part {
			msg := fmt.Sprintf("conflicting segment %q and %q produce same key", nd.Segment, part)
			return appErr.NewValidationError("path", msg)
		}

		if i == len(r.Parts)-1 {
			if nd.Root != "" && nd.Root != r.Path {
				msg := fmt.Sprintf("%q duplicate route detected: %q vs %q", r.Name, nd.Root, r.Path)
				return appErr.NewValidationError("path", msg)
			}
			nd.Root = r.Path
		}

		current = nd.Children
	}

	return nil
}

func (b *tree) normalizePart(r Route, part string) (string, string, error) {
	identPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	staticPattern := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

	if part == "*" {
		return b.cfg.Symbols.Wildcard, "wildcard", nil
	}

	if cut, found := strings.CutPrefix(part, ":"); found {
		if !identPattern.MatchString(cut) {
			msg := fmt.Sprintf("%q invalid param name: %s", r.Name, cut)

			return "", "", appErr.NewValidationError("path", msg)
		}
		return b.cfg.Symbols.Param, "param", nil
	}

	if !staticPattern.MatchString(part) {
		msg := fmt.Sprintf("%q invalid path segment: %s", r.Name, part)

		return "", "", appErr.NewValidationError("path", msg)
	}

	if strings.Contains(part, "-") {
		Parts := strings.Split(part, "-")
		for i := 1; i < len(Parts); i++ {
			if len(Parts[i]) > 0 {
				Parts[i] = strings.ToUpper(Parts[i][:1]) + Parts[i][1:]
			}
		}
		return strings.Join(Parts, ""), "static", nil
	}

	return part, "static", nil
}

func (b *tree) validateChild(r Route, current map[string]*Node, kind string) error {
	for _, c := range current {
		if kind == "wildcard" || c.Kind == "wildcard" {
			msg := fmt.Sprintf("%q path wildcard route cannot coexist with other routes at the same level", r.Name)

			return appErr.NewValidationError("path", msg)
		}

		if kind == "param" && c.Kind == "static" {
			msg := fmt.Sprintf("%q path param route conflicts with existing static route", r.Name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "static" && c.Kind == "param" {
			msg := fmt.Sprintf("%q path static route conflicts with existing param route", r.Name)

			return appErr.NewValidationError("path", msg)

		}

		if kind == "param" && c.Kind == "param" {
			msg := fmt.Sprintf("%q path multiple param routes at the same level are not allowed", r.Name)

			return appErr.NewValidationError("path", msg)
		}
	}

	return nil
}
