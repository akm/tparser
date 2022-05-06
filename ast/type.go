package ast

import (
	"fmt"

	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

// - TypeSection
//   ```
//   TYPE (TypeDecl ';')...
//   ```
type TypeSection []*TypeDecl // must implement InterfaceDecl

func (TypeSection) canBeInterfaceDecl() {}
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
	astcore.Decl
}

func (m *TypeDecl) Children() Nodes {
	return Nodes{m.Ident, m.Type}
}

func (m *TypeDecl) ToDeclarations() astcore.Declarations {
	return astcore.Declarations{astcore.NewDeclaration(m.Ident, m)}
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
type TypeId struct {
	Type
	UnitId *UnitId
	Ident  *Ident
	Ref    *astcore.Declaration // Actual Type object
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

func (*TypeId) isType() {}
func (m *TypeId) Children() Nodes {
	r := Nodes{}
	if m.UnitId != nil {
		r = append(r, m.UnitId)
	}
	r = append(r, m.Ident)
	return r
}
func (m *TypeId) String() string {
	if m.UnitId != nil {
		return fmt.Sprintf("%s.%s", m.UnitId.String(), m.Ident.Name)
	} else {
		return m.Ident.Name
	}
}
