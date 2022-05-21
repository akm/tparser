package ast

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

// - TypeSection
//   ```
//   TYPE (TypeDecl ';')...
//   ```
type TypeSection []*TypeDecl // must implement InterfaceDecl

func (TypeSection) canBeInterfaceDecl() {}
func (TypeSection) canBeDeclSection()   {}
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
	astcore.DeclNode
}

func (m *TypeDecl) Children() Nodes {
	return Nodes{m.Ident, m.Type}
}

func (m *TypeDecl) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
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

func (*TypeId) isType() {}
func (m *TypeId) Children() Nodes {
	r := Nodes{}
	if m.UnitId != nil {
		r = append(r, m.UnitId)
	}
	r = append(r, m.Ident)
	return r
}
