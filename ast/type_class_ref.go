package ast

type ClassRefType interface {
	Type
	IsClassRefType() bool
}

type CustomClassRefType struct {
	TypeId *TypeId
}

func NewCustomClassRefType(typeId *TypeId) *CustomClassRefType {
	return &CustomClassRefType{TypeId: typeId}
}

var _ ClassRefType = (*CustomClassRefType)(nil)

func (*CustomClassRefType) isType()              {}
func (*CustomClassRefType) IsClassRefType() bool { return true }
func (m *CustomClassRefType) Children() Nodes {
	return Nodes{m.TypeId}
}
