package astcore

import "github.com/pkg/errors"

// Ident with reference to Declaration
type IdentRef struct {
	*Ident
	Ref *Decl
	Node
}

func NewIdentRef(ident *Ident, ref *Decl) *IdentRef {
	return &IdentRef{Ident: ident, Ref: ref}
}

func (m *IdentRef) Children() Nodes {
	if m == nil {
		panic(errors.Errorf("IdentRef is nil. Something wrong"))
	}
	return Nodes{m.Ident}
}
