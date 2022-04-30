package ast

import (
	"strings"

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

func (*RealType) isType()       {}
func (*RealType) isSimpleType() {}

type RealType struct {
	Name Ident
}

func NewRealType(name interface{}) *RealType {
	switch v := name.(type) {
	case *RealType:
		return v
	case Ident:
		return &RealType{Name: v}
	default:
		panic(errors.Errorf("invalid type %T for NewRealType %+v", name, name))
	}
}

func (m *RealType) Children() Nodes {
	return Nodes{m.Name}
}

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

func (*OrdIdent) isType()        {}
func (*OrdIdent) isSimpleType()  {}
func (*OrdIdent) isOrdinalType() {}

type OrdIdent struct {
	Name Ident
}

func NewOrdIdent(name interface{}) *OrdIdent {
	switch v := name.(type) {
	case *OrdIdent:
		return v
	case Ident:
		return &OrdIdent{Name: v}
	case *Ident:
		return &OrdIdent{Name: *v}
	default:
		panic(errors.Errorf("invalid type %T for NewOrdIndent %+v", name, name))
	}
}

func (m *OrdIdent) Children() Nodes {
	return Nodes{m.Name}
}

// - EnumeratedType
//   ```
//   '(' EnumeratedTypeElement ','... ')'
//   ```
// - EnumeratedTypeElement
//   ```
//   Ident [ '=' ConstExpr ]
//   ```
func (EnumeratedType) isType()        {}
func (EnumeratedType) isSimpleType()  {}
func (EnumeratedType) isOrdinalType() {}

type EnumeratedType []*EnumeratedTypeElement

func (m EnumeratedType) Children() Nodes {
	r := make(Nodes, len(m))
	for i, e := range m {
		r[i] = e
	}
	return r
}

type EnumeratedTypeElement struct {
	Ident     Ident
	ConstExpr *ConstExpr
}

func (m *EnumeratedTypeElement) Children() Nodes {
	return Nodes{m.Ident, m.ConstExpr}.Compact()
}

// - SubrangeType
//   ```
//   ConstExpr '..' ConstExpr
//   ```
func (*SubrangeType) isType()        {}
func (*SubrangeType) isSimpleType()  {}
func (*SubrangeType) isOrdinalType() {}

type SubrangeType struct {
	Low  ConstExpr
	High ConstExpr
}

func (m *SubrangeType) Children() Nodes {
	return Nodes{&m.Low, &m.High}
}
