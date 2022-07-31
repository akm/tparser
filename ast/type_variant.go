package ast

// - VariantType
//   ```
//   VARIANT
//   ```
//   ```
//   OLEVARIANT
//   ```
type VariantType interface {
	Type
	IsVariantType() bool
}

func NewVariantType(v *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkVariantType, v.Name); decl != nil {
		return NewTypeId(v, decl)
	} else {
		return NewTypeId(v)
	}
}
