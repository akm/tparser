package ast

import (
	"github.com/akm/tparser/ast/astcore"
)

// - SimpleType
//   ```
//   (OrdinalType | RealType)
//   ```
type SimpleType interface {
	Type
	IsSimpleType() bool
}

// - RealType
//   ```
//   REAL48
//   ```
//   ```
//   REAL
//   ```
//   ```
//   SINGLE
//   ```
//   ```
//   DOUBLE
//   ```
//   ```
//   EXTENDED
//   ```
//   ```
//   CURRENCY
//   ```
//   ```
//   COMP
//   ```
type RealType interface {
	SimpleType
	IsRealType() bool
}

func NewRealType(ident *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkReal, ident.Name); decl != nil {
		return NewTypeId(ident, decl)
	} else {
		return NewTypeId(ident)
	}
}

// - OrdinalType
//   ```
//   (SubrangeType | EnumeratedType | OrdIdent)
//   ```

type OrdinalType interface {
	SimpleType
	IsOrdinalType() bool
}

// - OrdIdent
//   ```
//   SHORTINT
//   ```
//   ```
//   SMALLINT
//   ```
//   ```
//   INTEGER
//   ```
//   ```
//   BYTE
//   ```
//   ```
//   LONGINT
//   ```
//   ```
//   INT64
//   ```
//   ```
//   WORD
//   ```
//   ```
//   BOOLEAN
//   ```
//   ```
//   CHAR
//   ```
//   ```
//   WIDECHAR
//   ```
//   ```
//   LONGWORD
//   ```
//   ```
//   PCHAR
//   ```
type OrdIdent interface {
	OrdinalType
	IsOrdIdent() bool
}

func NewOrdIdent(ident *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkOrdIdent, ident.Name); decl != nil {
		return NewTypeId(ident, decl)
	} else {
		return NewTypeId(ident)
	}
}

// - EnumeratedType
//   ```
//   '(' EnumeratedTypeElement ','... ')'
//   ```
// - EnumeratedTypeElement
//   ```
//   Ident [ '=' ConstExpr ]
//   ```
type EnumeratedType []*EnumeratedTypeElement // must implement OrdinalType

var _ OrdinalType = (EnumeratedType)(nil)

func (EnumeratedType) isType()             {}
func (EnumeratedType) IsSimpleType() bool  { return true }
func (EnumeratedType) IsOrdinalType() bool { return true }
func (m EnumeratedType) Children() Nodes {
	r := make(Nodes, len(m))
	for i, e := range m {
		r[i] = e
	}
	return r
}

type EnumeratedTypeElement struct {
	astcore.DeclNode
	*Ident
	ConstExpr *ConstExpr
}

var _ Node = (*EnumeratedTypeElement)(nil)

func (m *EnumeratedTypeElement) Children() Nodes {
	r := Nodes{m.Ident}
	if m.ConstExpr != nil {
		r = append(r, m.ConstExpr)
	}
	return r
}
func (m *EnumeratedTypeElement) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

// - SubrangeType
//   ```
//   ConstExpr '..' ConstExpr
//   ```
type SubrangeType struct {
	OrdinalType
	Low  *ConstExpr
	High *ConstExpr
}

var _ OrdinalType = (*SubrangeType)(nil)

func (*SubrangeType) isType()             {}
func (*SubrangeType) IsSimpleType() bool  { return true }
func (*SubrangeType) IsOrdinalType() bool { return true }
func (m *SubrangeType) Children() Nodes   { return Nodes{m.Low, m.High} }
