package ast

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
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
	IsRealType() bool
	// implements
	SimpleType
}

func NewRealType(name interface{}) RealType {
	switch v := name.(type) {
	case RealType:
		return v
	case Ident:
		if decl := EmbeddedTypeDecl(EtkReal, v.Name); decl != nil {
			return NewTypeId(&v, decl)
		} else {
			return NewTypeId(&v)
		}
	case *Ident:
		if decl := EmbeddedTypeDecl(EtkReal, v.Name); decl != nil {
			return NewTypeId(v, decl)
		} else {
			return NewTypeId(v)
		}
	default:
		panic(errors.Errorf("invalid type %T for NewRealType %+v", name, name))
	}
}

// - OrdinalType
//   ```
//   (SubrangeType | EnumeratedType | OrdIdent)
//   ```

type OrdinalType interface {
	IsOrdinalType() bool
	// implements
	SimpleType
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
	IsOrdIdent() bool
	// implements
	OrdinalType
}

func NewOrdIdent(name interface{}) OrdIdent {
	switch v := name.(type) {
	case OrdIdent:
		return v
	case Ident:
		if decl := EmbeddedTypeDecl(EtkOrdIdent, v.Name); decl != nil {
			return NewTypeId(&v, decl)
		} else {
			return NewTypeId(&v)
		}
	case *Ident:
		if decl := EmbeddedTypeDecl(EtkOrdIdent, v.Name); decl != nil {
			return NewTypeId(v, decl)
		} else {
			return NewTypeId(v)
		}
	default:
		panic(errors.Errorf("invalid type %T for NewOrdIndent %+v", name, name))
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

func (*SubrangeType) isType()             {}
func (*SubrangeType) IsSimpleType() bool  { return true }
func (*SubrangeType) IsOrdinalType() bool { return true }
func (m *SubrangeType) Children() Nodes   { return Nodes{m.Low, m.High} }
