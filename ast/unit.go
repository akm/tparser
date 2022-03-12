package ast

func (*Unit) isGoal() {}

type Unit struct {
	Ident                 Ident
	PortabilityDirective  *PortabilityDirective // optional
	InterfaceSection      *InterfaceSection
	ImplementationSection *ImplementationSection
	InitSection           *InitSection // optional
}

type InterfaceSection struct {
	UsesClause *UsesClause // optional
}

type InterfaceDecl interface {
	canBeInterfaceDecl()
}

type ImplementationSection struct {
}

type InitSection struct {
}

// UnitId: <unit-identifier>
type UnitId Ident
