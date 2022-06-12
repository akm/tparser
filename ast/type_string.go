package ast

import (
	"github.com/pkg/errors"
)

// - StringType
//   ```
//   STRING
//   ```
//   ```
//   ANSISTRING
//   ```
//   ```
//   WIDESTRING
//   ```
//   ```
//   STRING '[' ConstExpr ']'
//   ```
type StringType interface {
	IsStringType() bool
	// implements
	Type
}

func NewStringType(name interface{}) StringType {
	switch v := name.(type) {
	case StringType:
		return v
	case *Ident:
		if decl := EmbeddedTypeDecl(EtkStringType, v.Name); decl != nil {
			return NewTypeId(v, decl)
		} else {
			return NewTypeId(v)
		}
	default:
		panic(errors.Errorf("invalid type %T for NewStringType %+v", name, name))
	}
}

type FixedStringType struct {
	StringType
	Length *ConstExpr
}

func NewFixedStringType(name interface{}, length *ConstExpr) *FixedStringType {
	return &FixedStringType{StringType: NewStringType(name), Length: length}
}

func (*FixedStringType) isType()            {}
func (*FixedStringType) IsStringType() bool { return true }
func (m *FixedStringType) Children() Nodes {
	return Nodes{m.StringType, m.Length}
}
