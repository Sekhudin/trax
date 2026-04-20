package ts

import (
	"fmt"
	"strings"
)

func Line(sl ...string) string {
	return strings.Join(sl, "\n")
}

func ToTypeLiteral(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func ToStringUnion(s string) string {
	return fmt.Sprintf("%q", s)
}

func OrganizeImport(mods []string) string {
	m := make(map[string]struct{})
	r := make([]string, 0, len(mods))

	for _, mod := range mods {
		if _, oke := m[mod]; oke {
			continue
		}
		m[mod] = struct{}{}
		r = append(r, mod)
	}

	return strings.Join(r, ", ")
}

func OrganizeStringUnion(ss []string) string {
	if len(ss) == 1 {
		return ToStringUnion(ss[0])
	}

	m := make(map[string]struct{})
	r := make([]string, 0, len(ss))

	for _, s := range ss {
		if _, oke := m[s]; oke {
			continue
		}
		m[s] = struct{}{}
		r = append(r, ToStringUnion(s))
	}

	return strings.Join(r, " | ")
}
