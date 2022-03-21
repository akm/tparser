package ast

func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

type ConstantDecl struct {
	Ident                Ident
	Type                 Type
	ConstExpr            ConstExpr
	PortabilityDirective *PortabilityDirective
}

type ConstExpr = Expression
