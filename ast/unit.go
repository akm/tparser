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
	UsesClause     *UsesClause // optional
	InterfaceDecls []InterfaceDecl
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

func NewUnitId(name string) *UnitId {
	r := UnitId(name)
	return &r
}

func (u *UnitId) String() string {
	if u == nil {
		return ""
	} else {
		return string(*u)
	}
}

type QualId struct {
	UnitId *UnitId
	Ident  Ident
}
