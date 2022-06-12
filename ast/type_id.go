package ast

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

// - TypeId
//   ```
//   [UnitId '.'] <type-identifier>
//   ```
type TypeId struct {
	Type
	UnitId *UnitId
	Ident  *Ident
	Ref    *astcore.Decl // Actual Type object
}

func NewTypeId(unitIdOrIdent interface{}, args ...interface{}) *TypeId {
	var ref *astcore.Decl
	if len(args) > 0 {
		if v, ok := args[len(args)-1].(*astcore.Decl); ok {
			ref = v
			args = args[:len(args)-1]
		}
	}
	if len(args) == 0 {
		return &TypeId{Ident: NewIdentFrom(unitIdOrIdent), Ref: ref}
	} else if len(args) == 1 {
		return &TypeId{
			UnitId: NewUnitId(unitIdOrIdent),
			Ident:  NewIdentFrom(args[0]),
			Ref:    ref,
		}
	} else {
		panic(errors.Errorf("too many arguments for NewTypeId: %v, %v", unitIdOrIdent, args))
	}
}

func (m *TypeId) Children() Nodes {
	r := Nodes{}
	if m.UnitId != nil {
		r = append(r, m.UnitId)
	}
	r = append(r, m.Ident)
	return r
}

func (*TypeId) isType() {}
func (m *TypeId) IsSimpleType() bool {
	if m.Ref == nil {
		return false
	}
	if m.Ref.Node == nil {
		return false
	}
	decl, ok := m.Ref.Node.(*TypeDecl)
	if !ok {
		return false
	}
	simpleType, ok := decl.Type.(SimpleType)
	if !ok {
		return false
	}
	return simpleType.IsSimpleType()
}

func (m *TypeId) IsOrdinalType() bool {
	if m.Ref == nil {
		return false
	}
	if m.Ref.Node == nil {
		return false
	}
	decl, ok := m.Ref.Node.(*TypeDecl)
	if !ok {
		return false
	}
	ordinalType, ok := decl.Type.(OrdinalType)
	if !ok {
		return false
	}
	return ordinalType.IsOrdinalType()
}

func (m *TypeId) IsRealType() bool {
	if m.Ref == nil {
		return false
	}
	if m.Ref.Node == nil {
		return false
	}
	decl, ok := m.Ref.Node.(*TypeDecl)
	if !ok {
		return false
	}
	ordinalType, ok := decl.Type.(RealType)
	if !ok {
		return false
	}
	return ordinalType.IsRealType()
}

func (m *TypeId) IsOrdIdent() bool {
	if m.Ref == nil {
		return false
	}
	if m.Ref.Node == nil {
		return false
	}
	decl, ok := m.Ref.Node.(*TypeDecl)
	if !ok {
		return false
	}
	ordIdent, ok := decl.Type.(OrdIdent)
	if !ok {
		return false
	}
	return ordIdent.IsOrdIdent()
}
