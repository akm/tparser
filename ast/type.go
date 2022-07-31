package ast

import (
	"github.com/akm/tparser/ast/astcore"
)

// - TypeSection
//   ```
//   TYPE (TypeDecl ';')...
//   ```
type TypeSection []*TypeDecl

var _ Node = (TypeSection)(nil)
var _ InterfaceDecl = (TypeSection)(nil)

func (TypeSection) canBeInterfaceDecl() {}
func (TypeSection) canBeDeclSection()   {}
func (s TypeSection) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}
func (s TypeSection) GetDeclNodes() astcore.DeclNodes {
	r := make(astcore.DeclNodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - TypeDecl
//   ```
//   Ident '=' [TYPE] Type [PortabilityDirective]
//   ```
//   ```
//   Ident '=' [TYPE] RestrictedType [PortabilityDirective]
//   ```
type TypeDecl struct {
	*Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
}

var _ astcore.DeclNode = (*TypeDecl)(nil)

func (m *TypeDecl) Children() Nodes {
	return Nodes{m.Ident, m.Type}
}

func (m *TypeDecl) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

// - Type
//   ```
//   TypeId
//   ```
//   ```
//   SimpleType
//   ```
//   ```
//   StrucType
//   ```
//   ```
//   PointerType
//   ```
//   ```
//   StringType
//   ```
//   ```
//   ProcedureType
//   ```
//   ```
//   VariantType
//   ```
//   ```
//   ClassRefType
//   ```
type Type interface {
	Node
	isType()
}
