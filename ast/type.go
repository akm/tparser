package ast

import (
	"strings"

	"github.com/akm/tparser/ext"
)

func (TypeSection) canBeInterfaceDecl() {}

type TypeSection []*TypeDecl

type TypeDecl struct {
	Ident                Ident
	Type                 Type
	PortabilityDirective *PortabilityDirective
}

type Type interface {
	isRestrictedType() bool
}

func (*TypeId) isRestrictedType() bool { return false }

// TypeId: [UnitId '.'] <type-identifier>
type TypeId struct {
	UnitId *UnitId
	Ident  Ident
}

func IsRealTypeName(w string) bool {
	return realTypeNames.Include(strings.ToUpper(w))
}

var realTypeNames = ext.Strings{
	"REAL48",
	"REAL",
	"SINGLE",
	"DOUBLE",
	"EXTENDED",
	"CURRENCY",
	"COMP",
}.Set()

func (*RealType) isRestrictedType() bool { return false }

type RealType struct {
	Name Ident
}

func IsOrdIdentName(w string) bool {
	return ordIdentNames.Include(strings.ToUpper(w))
}

var ordIdentNames = ext.Strings{
	// Integer types
	"INTEGER",
	"CARDINAL",
	"SHORTINT",
	"SMALLINT",
	"LONGINT",
	"INT64",
	"BYTE",
	"WORD",
	"LONGWORD",

	// Character types
	"CHAR",
	"ANSICHAR",
	"WIDECHAR",

	// Boolean types
	"BOOLEAN",

	// The following are in String Type
	// "PCHAR",
	// "PANSICHAR",
	// "PWIDECHAR",
}.Set()

func (*OrdIdent) isRestrictedType() bool { return false }

type OrdIdent struct {
	Name Ident
}
