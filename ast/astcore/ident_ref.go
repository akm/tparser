package astcore

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
	return Nodes{m.Ident}
}
