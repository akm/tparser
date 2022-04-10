package ast

func (TypeSection) canBeInterfaceDecl() {}

type TypeSection []*TypeDecl

type TypeDecl struct {
	Ident                Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
}

type Type interface {
	isType()
}

// TypeId: [UnitId '.'] <type-identifier>
func (*TypeId) isType() {}

type TypeId struct {
	UnitId *UnitId
	Ident  Ident
}
