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
}

type InterfaceDecl interface {
	canBeInterfaceDecl()
}

type ImplementationSection struct {
}

type InitSection struct {
}