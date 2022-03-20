package ast

func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

type ConstantDecl struct {
	Ident                Ident
	Type                 Type
	ConstExpr            ConstExpr
	PortabilityDirective *PortabilityDirective
}

// TODO implement ConstExpr
type ConstExpr struct {
	Value string
}
