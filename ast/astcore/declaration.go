package astcore

type Declaration struct {
	*Ident
	Node Node
}

func NewDeclaration(ident *Ident, node Node) *Declaration {
	return &Declaration{Ident: ident, Node: node}
}

type Declarations []*Declaration

func NewDeclarations(identList IdentList, node Node) Declarations {
	r := make(Declarations, len(identList))
	for idx, i := range identList {
		r[idx] = NewDeclaration(i, node)
	}
	return r
}
