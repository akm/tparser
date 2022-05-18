package astcore

// Ident with reference to Declaration
type IdentRef struct {
	*Ident
	Ref *Declaration
	Node
}

func NewIdentRef(ident *Ident, ref *Declaration) *IdentRef {
	return &IdentRef{Ident: ident, Ref: ref}
}

func (m *IdentRef) Children() Nodes {
	return Nodes{m.Ident}
}
