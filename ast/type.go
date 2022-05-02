package ast

import "github.com/pkg/errors"

// - TypeSection
//   ```
//   TYPE (TypeDecl ';')...
//   ```
func (TypeSection) canBeInterfaceDecl() {}

type TypeSection []*TypeDecl

func (s TypeSection) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - TypeDecl
//   ```
//   Ident '=' [TYPE] Type [PortabilityDirective]
//   ```
//   ```
//   Ident '=' [TYPE] RestrictedType [PortabilityDirective]
//   ```
type TypeDecl struct {
	*Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
}

func (m *TypeDecl) Children() Nodes {
	return Nodes{m.Ident, m.Type}
}

// - Type
//   ```
//   TypeId
//   ```
//   ```
//   SimpleType
//   ```
//   ```
//   StrucType
//   ```
//   ```
//   PointerType
//   ```
//   ```
//   StringType
//   ```
//   ```
//   ProcedureType
//   ```
//   ```
//   VariantType
//   ```
//   ```
//   ClassRefType
//   ```
type Type interface {
	Node
	isType()
}

// - TypeId
//   ```
//   [UnitId '.'] <type-identifier>
//   ```
func (*TypeId) isType() {}

type TypeId struct {
	UnitId *UnitId
	Ident  *Ident
}

func NewTypeId(unitIdOrIdent interface{}, args ...interface{}) *TypeId {
	if len(args) == 0 {
		return &TypeId{Ident: NewIdentFrom(unitIdOrIdent)}
	} else if len(args) == 1 {
		return &TypeId{
			UnitId: NewUnitId(unitIdOrIdent),
			Ident:  NewIdentFrom(args[0]),
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
