package ast

func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

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
