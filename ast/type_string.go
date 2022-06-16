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
	IsStringType() bool
	// implements
	Type
}

func NewStringType(ident *Ident) StringType {
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

func NewFixedStringType(ident *Ident, length *ConstExpr) *FixedStringType {
	return &FixedStringType{StringType: NewStringType(ident), Length: length}
}

func (*FixedStringType) isType()            {}
func (*FixedStringType) IsStringType() bool { return true }
func (m *FixedStringType) Children() Nodes {
	return Nodes{m.StringType, m.Length}
}
