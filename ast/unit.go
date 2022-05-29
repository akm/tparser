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
	DeclMap               astcore.DeclMap

	Namespace
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

func (m *Unit) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

func (m *Unit) GetIdent() *Ident {
	return m.Ident
}

func (m *Unit) GetDeclMap() astcore.DeclMap {
	return m.DeclMap
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

func (s Units) DeclMaps() astcore.DeclMaps {
	r := make(astcore.DeclMaps, len(s))
	for i, m := range s {
		r[i] = m.DeclMap
	}
	return r
}

func (s Units) Compact() Units {
	r := Units{}
	for _, i := range s {
		if i != nil {
			r = append(r, i)
		}
	}
	return r
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
	GetDeclNodes() astcore.DeclNodes
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
//   [NamespaceId '.'] Ident
//   ```
//
// - NamespaceId
//   ```
//   <unit-identifier>
//   ```
//   ```
//   <program-identifier>
//   ```
type QualId struct {
	NamespaceId *IdentRef
	Ident       *IdentRef
}

func NewQualId(unitId *IdentRef, ident *IdentRef) *QualId {
	return &QualId{NamespaceId: unitId, Ident: ident}
}

func (m *QualId) Children() Nodes {
	r := Nodes{}
	if m.NamespaceId != nil {
		r = append(r, m.NamespaceId)
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
