package ast

import "github.com/pkg/errors"

// - UsesClause
//   ```
//   USES IdentList ';'
//   ```
// Actually USES has not only IdentList.
type UsesClause []*UnitRef

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

type UnitRef struct {
	*Ident
	Path *string
}

func NewUnitRef(name interface{}, paths ...string) *UnitRef {
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
	r := &UnitRef{Ident: nameIdent}
	if len(paths) > 1 {
		panic(errors.Errorf("too many paths: %v for NewUnitPath", paths))
	}
	if len(paths) == 1 {
		s := paths[0]
		r.Path = &s
	}
	return r
}

func (m *UnitRef) Children() Nodes {
	return Nodes{m.Ident}
}
