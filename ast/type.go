package ast

import "github.com/pkg/errors"

// - TypeSection
//   ```
//   TYPE (TypeDecl ';')...
//   ```
func (TypeSection) canBeInterfaceDecl() {}

type TypeSection []*TypeDecl

// - TypeDecl
//   ```
//   Ident '=' [TYPE] Type [PortabilityDirective]
//   ```
//   ```
//   Ident '=' [TYPE] RestrictedType [PortabilityDirective]
//   ```
type TypeDecl struct {
	Ident                Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
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
	isType()
}

// - TypeId
//   ```
//   [UnitId '.'] <type-identifier>
//   ```
func (*TypeId) isType() {}

type TypeId struct {
	UnitId *UnitId
	Ident  Ident
}

func NewTypeId(unitIdOrIdent interface{}, args ...interface{}) *TypeId {
	if len(args) == 0 {
		return &TypeId{Ident: NewIdent(unitIdOrIdent)}
	} else if len(args) == 1 {
		return &TypeId{
			UnitId: NewUnitId(unitIdOrIdent),
			Ident:  NewIdent(args[0]),
		}
	} else {
		panic(errors.Errorf("too many arguments for NewTypeId: %v, %v", unitIdOrIdent, args))
	}
}
