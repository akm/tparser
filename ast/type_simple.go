package ast

import (
	"strings"

	"github.com/akm/tparser/ext"
)

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

type RealType struct {
	Name Ident
}

type OrdinalType interface {
	isOrdinalType()
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

func (*OrdIdent) isOrdinalType() {}

type OrdIdent struct {
	Name Ident
}

func (EnumeratedType) isOrdinalType() {}

type EnumeratedType []*EnumeratedTypeElement

type EnumeratedTypeElement struct {
	Ident     Ident
	ConstExpr *ConstExpr
}

func (*SubrangeType) isOrdinalType() {}

type SubrangeType struct {
	Low  ConstExpr
	High ConstExpr
}
