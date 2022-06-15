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

func (m *TypeId) getRefNodeDecl() *TypeDecl {
	if m.Ref == nil {
		return nil
	}
	if m.Ref.Node == nil {
		return nil
	}
	decl, ok := m.Ref.Node.(*TypeDecl)
	if !ok {
		return nil
	}
	return decl
}

func (*TypeId) isType() {}
func (m *TypeId) IsSimpleType() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	simpleType, ok := decl.Type.(SimpleType)
	if !ok {
		return false
	}
	return simpleType.IsSimpleType()
}

func (m *TypeId) IsOrdinalType() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	ordinalType, ok := decl.Type.(OrdinalType)
	if !ok {
		return false
	}
	return ordinalType.IsOrdinalType()
}

func (m *TypeId) IsRealType() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	ordinalType, ok := decl.Type.(RealType)
	if !ok {
		return false
	}
	return ordinalType.IsRealType()
}

func (m *TypeId) IsOrdIdent() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	ordIdent, ok := decl.Type.(OrdIdent)
	if !ok {
		return false
	}
	return ordIdent.IsOrdIdent()
}

func (m *TypeId) IsStringType() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	ordIdent, ok := decl.Type.(StringType)
	if !ok {
		return false
	}
	return ordIdent.IsStringType()
}

func (m *TypeId) IsPointerType() bool {
	decl := m.getRefNodeDecl()
	if decl == nil {
		return false
	}
	ordIdent, ok := decl.Type.(PointerType)
	if !ok {
		return false
	}
	return ordIdent.IsPointerType()
}
