package ast

// - PointerType
//   ```
//   '^' TypeId [PortabilityDirective]
//   ```
type PointerType interface {
	Type
	IsPointerType() bool
}

func NewEmbeddedPointerType(v *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkPointerType, v.Name); decl != nil {
		return NewTypeId(v, decl)
	} else {
		return NewTypeId(v)
	}
}

func NewPointerType(ident *Ident) *TypeId {
	if decl := EmbeddedTypeDecl(EtkPointerType, ident.Name); decl != nil {
		return NewTypeId(ident, decl)
	} else {
		return NewTypeId(ident)
	}
}

type CustomPointerType struct {
	TypeId *TypeId
}

var _ PointerType = (*CustomPointerType)(nil)

func NewCustomPointerType(typeId *TypeId) *CustomPointerType {
	return &CustomPointerType{TypeId: typeId}
}

func (*CustomPointerType) isType()             {}
func (*CustomPointerType) IsPointerType() bool { return true }
func (m *CustomPointerType) Children() Nodes {
	return Nodes{m.TypeId}
}
