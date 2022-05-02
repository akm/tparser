package ast

import "strings"

func (*Unit) isGoal() {}

// - Unit
//   ```
//   UNIT Ident [PortabilityDirective] ';'
//   InterfaceSection
//   ImplementationSection
//   InitSection '.'
//   ```
type Unit struct {
	Path                  string
	Ident                 *Ident
	PortabilityDirective  *PortabilityDirective // optional
	InterfaceSection      *InterfaceSection
	ImplementationSection *ImplementationSection
	InitSection           *InitSection // optional
}

func (m *Unit) GetPath() string {
	return m.Path
}

func (m *Unit) Children() Nodes {
	r := Nodes{m.Ident}
	if m.InterfaceSection != nil {
		r = append(r, m.InterfaceSection)
	}
	if m.ImplementationSection != nil {
		r = append(r, m.ImplementationSection)
	}
	if m.InitSection != nil {
		r = append(r, m.InitSection)
	}
	return r
}

type Units []*Unit

func (s Units) ByName(name string) *Unit {
	key := strings.ToLower(name)
	for _, u := range s {
		if strings.ToLower(u.Ident.Name) == key {
			return u
		}
	}
	return nil
}

// - InterfaceSection
//   ```
//   INTERFACE
//   [UsesClause]
//   [InterfaceDecl]...
//   ```
type InterfaceSection struct {
	UsesClause     UsesClause // optional
	InterfaceDecls InterfaceDecls
}

func (m *InterfaceSection) Children() Nodes {
	r := Nodes{}
	if m.UsesClause != nil {
		r = append(r, m.UsesClause)
	}
	if m.InterfaceDecls != nil {
		r = append(r, m.InterfaceDecls)
	}
	return r
}

// - InterfaceDecl
//   ```
//   ConstSection
//   ```
//   ```
//   TypeSection
//   ```
//   ```
//   VarSection
//   ```
//   ```
//   ExportedHeading
//   ```
type InterfaceDecl interface {
	Node
	canBeInterfaceDecl()
}

type InterfaceDecls []InterfaceDecl

func (s InterfaceDecls) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ImplementationSection
//   ```
//   IMPLEMENTATION
//   [UsesClause]
//   [DeclSection]...
//   [ExportsStmt]...
//   ```
type ImplementationSection struct {
}

func (m *ImplementationSection) Children() Nodes {
	return Nodes{}
}

// - InitSection
//   ```
//   INITIALIZATION StmtList [FINALIZATION StmtList] END
//   ```
//   ```
//   BEGIN StmtList END
//   ```
//   ```
//   END
//   ```
type InitSection struct {
}

func (m *InitSection) Children() Nodes {
	return Nodes{}
}

// - UnitId
//   ```
//   <unit-identifier>
//   ```
type UnitId Ident

func NewUnitId(name interface{}) *UnitId {
	r := UnitId(*NewIdentFrom(name))
	return &r
}

func (u *UnitId) String() string {
	if u == nil {
		return ""
	} else {
		return u.Name
	}
}

func (m *UnitId) Children() Nodes {
	return Nodes{(*Ident)(m)}
}

// - QualId
//   ```
//   [UnitId '.'] Ident
//   ```
type QualId struct {
	UnitId *UnitId
	Ident  *Ident
}

func (m *QualId) Children() Nodes {
	r := Nodes{}
	if m.UnitId != nil {
		r = append(r, m.UnitId)
	}
	r = append(r, m.Ident)
	return r
}
