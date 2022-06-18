package ast

// - VariantType
//   ```
//   VARIANT
//   ```
//   ```
//   OLEVARIANT
//   ```
type VariantType interface {
	IsVariantType() bool
	// implements
	Type
}

func NewVariantType(v *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkVariantType, v.Name); decl != nil {
		return NewTypeId(v, decl)
	} else {
		return NewTypeId(v)
	}
}
