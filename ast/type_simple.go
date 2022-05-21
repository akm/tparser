package ast

import (
	"strings"

	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/pkg/errors"
)

// - SimpleType
//   ```
//   (OrdinalType | RealType)
//   ```
type SimpleType interface {
	Type
	isSimpleType()
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
func IsRealTypeName(w string) bool {
	return realTypeNames.Include(strings.ToUpper(w))
}

var realTypeNames = ext.Strings{
	"REAL48",
	"REAL",
	"SINGLE",
	"DOUBLE",
	"EXTENDED",
	"CURRENCY",
	"COMP",
}.Set()

type RealType struct {
	SimpleType
	*Ident
}

func NewRealType(name interface{}) *RealType {
	switch v := name.(type) {
	case *RealType:
		return v
	case Ident:
		return &RealType{Ident: &v}
	case *Ident:
		return &RealType{Ident: v}
	default:
		panic(errors.Errorf("invalid type %T for NewRealType %+v", name, name))
	}
}

func (*RealType) isType()           {}
func (*RealType) isSimpleType()     {}
func (m *RealType) Children() Nodes { return Nodes{m.Ident} }

// - OrdinalType
//   ```
//   (SubrangeType | EnumeratedType | OrdIdent)
//   ```

type OrdinalType interface {
	SimpleType
	isOrdinalType()
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
func IsOrdIdentName(w string) bool {
	return ordIdentNames.Include(strings.ToUpper(w))
}

var ordIdentNames = ext.Strings{
	// Integer types
	"INTEGER",
	"CARDINAL",
	"SHORTINT",
	"SMALLINT",
	"LONGINT",
	"INT64",
	"BYTE",
	"WORD",
	"LONGWORD",

	// Character types
	"CHAR",
	"ANSICHAR",
	"WIDECHAR",

	// Boolean types
	"BOOLEAN",

	// The following are in String Type
	// "PCHAR",
	// "PANSICHAR",
	// "PWIDECHAR",
}.Set()

type OrdIdent struct {
	OrdinalType
	*Ident
}

func NewOrdIdent(name interface{}) *OrdIdent {
	switch v := name.(type) {
	case *OrdIdent:
		return v
	case Ident:
		return &OrdIdent{Ident: &v}
	case *Ident:
		return &OrdIdent{Ident: v}
	default:
		panic(errors.Errorf("invalid type %T for NewOrdIndent %+v", name, name))
	}
}

func (*OrdIdent) isType()           {}
func (*OrdIdent) isSimpleType()     {}
func (*OrdIdent) isOrdinalType()    {}
func (m *OrdIdent) Children() Nodes { return Nodes{m.Ident} }

// - EnumeratedType
//   ```
//   '(' EnumeratedTypeElement ','... ')'
//   ```
// - EnumeratedTypeElement
//   ```
//   Ident [ '=' ConstExpr ]
//   ```
type EnumeratedType []*EnumeratedTypeElement // must implement OrdinalType

func (EnumeratedType) isType()        {}
func (EnumeratedType) isSimpleType()  {}
func (EnumeratedType) isOrdinalType() {}
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

func (*SubrangeType) isType()           {}
func (*SubrangeType) isSimpleType()     {}
func (*SubrangeType) isOrdinalType()    {}
func (m *SubrangeType) Children() Nodes { return Nodes{m.Low, m.High} }
