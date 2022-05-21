package astcore

type Decl struct {
	*Ident
	Node Node
}

func NewDeclaration(ident *Ident, node Node) *Decl {
	return &Decl{Ident: ident, Node: node}
}

type Decls []*Decl

func NewDeclarations(identList IdentList, node Node) Decls {
	r := make(Decls, len(identList))
	for idx, i := range identList {
		r[idx] = NewDeclaration(i, node)
	}
	return r
}
