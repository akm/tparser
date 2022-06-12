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
