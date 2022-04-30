package ast

func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

func (s ConstSection) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

type ConstantDecl struct {
	Ident                Ident
	Type                 Type
	ConstExpr            ConstExpr
	PortabilityDirective *PortabilityDirective
}

func (m *ConstantDecl) Children() Nodes {
	r := Nodes{&m.Ident}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	r = append(r, &m.ConstExpr)
	return r
}

type ConstExpr = Expression

func NewConstExpr(arg interface{}) *ConstExpr {
	return NewExpression(arg)
}
