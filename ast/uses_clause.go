package ast

import (
	"path/filepath"
	"strings"

	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

// - UsesClause
//   ```
//   USES IdentList ';'
//   ```
// In the uses clause of a program or library, any unit name may be followed
// by the reserved word in and the name of a source file, with or without a
// directory path, in single quotation marks; directory paths can be absolute
// or relative. Examples:
//     uses Windows, Messages, SysUtils, Strings in 'C:\Classes\Strings.pas', Classes;
type UsesClause []*UsesClauseItem

func (s UsesClause) IdentList() IdentList {
	var ids IdentList
	for _, u := range s {
		ids = append(ids, u.Ident)
	}
	return ids
}

func (s UsesClause) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

func (s UsesClause) Find(name string) *UsesClauseItem {
	k := strings.ToLower(name)
	for _, u := range s {
		if strings.ToLower(u.Ident.Name) == k {
			return u
		}
	}
	return nil
}

func (s UsesClause) Units() Units {
	r := make(Units, len(s))
	for i, m := range s {
		r[i] = m.Unit
	}
	return r
}

type UsesClauseItem struct {
	*Ident
	Path *string
	Unit *Unit
	astcore.DeclNode
}

func NewUnitRef(name interface{}, paths ...string) *UsesClauseItem {
	var nameIdent *Ident
	switch v := name.(type) {
	case Ident:
		nameIdent = &v
	case *Ident:
		nameIdent = v
	case string:
		nameIdent = NewIdentFrom(v)
	default:
		panic(errors.Errorf("invalid type %T", name))
	}
	r := &UsesClauseItem{Ident: nameIdent}
	if len(paths) > 1 {
		panic(errors.Errorf("too many paths: %v for NewUnitPath", paths))
	}
	if len(paths) == 1 {
		s := paths[0]
		r.Path = &s
	}
	return r
}

func (m *UsesClauseItem) Children() Nodes {
	return Nodes{m.Ident}
}

func (m *UsesClauseItem) UnquotedPath() string {
	if m.Path == nil {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(*m.Path, "'"), "'")
}

func (m *UsesClauseItem) EffectivePath() string {
	origPath := m.UnquotedPath()
	return strings.ReplaceAll(origPath, "\\", string([]rune{filepath.Separator}))
}

func (m *UsesClauseItem) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}
