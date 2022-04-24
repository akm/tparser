package ast

func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

func (s ConstSection) Children() []Node {
	r := make([]Node, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

type ConstantDecl struct {
	CodeBlockNode
	Ident                Ident
	Type                 Type
	ConstExpr            ConstExpr
	PortabilityDirective *PortabilityDirective
}

type ConstExpr = Expression

func NewConstExpr(arg interface{}) *ConstExpr {
	return NewExpression(arg)
}
