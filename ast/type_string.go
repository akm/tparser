package ast

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
	Type
	IsStringType() bool
}

func NewStringType(ident *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkStringType, ident.Name); decl != nil {
		return NewTypeId(ident, decl)
	} else {
		return NewTypeId(ident)
	}
}

type FixedStringType struct {
	StringType
	Length *ConstExpr
}

var _ StringType = (*FixedStringType)(nil)

func NewFixedStringType(ident *Ident, length *ConstExpr) *FixedStringType {
	return &FixedStringType{StringType: NewStringType(ident), Length: length}
}

func (*FixedStringType) isType()            {}
func (*FixedStringType) IsStringType() bool { return true }
func (m *FixedStringType) Children() Nodes {
	return Nodes{m.StringType, m.Length}
}
