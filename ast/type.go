package ast

func (TypeSection) canBeInterfaceDecl() {}

type TypeSection []*TypeDecl

type TypeDecl struct {
	Ident                Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
}

type Type interface {
}

// TypeId: [UnitId '.'] <type-identifier>
type TypeId struct {
	UnitId *UnitId
	Ident  Ident
}
