package ast

// - PointerType
//   ```
//   '^' TypeId [PortabilityDirective]
//   ```
type PointerType interface {
	IsPointerType() bool
	// implements
	Type
}

func NewEmbeddedPointerType(v *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkPointerType, v.Name); decl != nil {
		return NewTypeId(v, decl)
	} else {
		return NewTypeId(v)
	}
}

type CustomPointerType struct {
	TypeId *TypeId
	// implements
	PointerType
}

func NewCustomPointerType(typeId *TypeId) *CustomPointerType {
	return &CustomPointerType{TypeId: typeId}
}

func (*CustomPointerType) isType()             {}
func (*CustomPointerType) IsPointerType() bool { return true }
func (m *CustomPointerType) Children() Nodes {
	return Nodes{m.TypeId}
}
