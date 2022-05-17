package ast

import (
	"strings"

	"github.com/akm/tparser/ast/astcore"
)

// - Unit
//   ```
//   UNIT Ident [PortabilityDirective] ';'
//   InterfaceSection
//   ImplementationSection
//   [InitSection] '.'
//   ```
type Unit struct {
	Path string
	*Ident
	PortabilityDirective  *PortabilityDirective // optional
	InterfaceSection      *InterfaceSection
	ImplementationSection *ImplementationSection
	InitSection           *InitSection // optional
	DeclarationMap        astcore.DeclarationMap
	Goal
}

func (*Unit) isGoal() {}
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

func (m *Unit) ToDeclarations() astcore.Declarations {
	return astcore.Declarations{astcore.NewDeclaration(m.Ident, m)}
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
	UsesClause   UsesClause // optional
	DeclSections DeclSections
	ExportsStmts ExportsStmts
	Node
}

func (m *ImplementationSection) Children() Nodes {
	r := Nodes{}
	if m.UsesClause != nil {
		r = append(r, m.UsesClause)
	}
	if m.DeclSections != nil {
		r = append(r, m.DeclSections)
	}
	if m.ExportsStmts != nil {
		r = append(r, m.ExportsStmts)
	}
	return r
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
	InitializationStmts StmtList
	FinalizationStmts   StmtList
}

func (m *InitSection) Children() Nodes {
	r := Nodes{m.InitializationStmts}
	if m.FinalizationStmts != nil {
		r = append(r, m.FinalizationStmts)
	}
	return r
}

// - UnitId
//   ```
//   <unit-identifier>
//   ```
type UnitId = Ident

func NewUnitId(name interface{}) *UnitId {
	switch v := name.(type) {
	case *UnitId:
		return v
	default:
		return NewIdentFrom(name)
	}
}

// - QualId
//   ```
//   [UnitId '.'] Ident
//   ```
type QualId struct {
	UnitId *UnitId
	Ident  *Ident
	Ref    *astcore.Declaration // Actual Type object
}

func NewQualId(unitId *UnitId, ident *Ident, refs ...*astcore.Declaration) *QualId {
	r := &QualId{UnitId: unitId, Ident: ident}
	if len(refs) > 0 {
		r.Ref = refs[0]
	}
	return r
}

func (m *QualId) Children() Nodes {
	r := Nodes{}
	if m.UnitId != nil {
		r = append(r, m.UnitId)
	}
	r = append(r, m.Ident)
	return r
}

// QualIds
type QualIds []*QualId // must implement Node

func (s QualIds) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}
