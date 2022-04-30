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
	return Nodes{&m.Ident, m.Type, &m.ConstExpr}.Compact()
}

type ConstExpr = Expression

func NewConstExpr(arg interface{}) *ConstExpr {
	return NewExpression(arg)
}
