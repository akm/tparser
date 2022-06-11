package ast

import (
	"strings"

	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

type EmbeddedTypeKind int

const (
	EtkReal EmbeddedTypeKind = iota + 1
	EtkOrdIdent
	EtkString
)

type TypeEmbedded struct {
	Kind  EmbeddedTypeKind
	Ident *Ident
	// implements
	Type
	SimpleType
	OrdinalType
}

func newTypeEmbedded(kind EmbeddedTypeKind, name string) *TypeEmbedded {
	return &TypeEmbedded{
		Kind: kind,
		Ident: &astcore.Ident{
			Name:     name,
			Location: &astcore.Location{Path: "embedded"},
		},
	}
}

func (m *TypeEmbedded) Children() Nodes {
	return Nodes{m.Ident}
}

func (*TypeEmbedded) isType() {}
func (m *TypeEmbedded) IsSimpleType() bool {
	return m.IsRealType() || m.IsOrdIdent()
}
func (m *TypeEmbedded) IsRealType() bool {
	return m.Kind == EtkReal
}

func (m *TypeEmbedded) IsOrdinalType() bool {
	return m.IsOrdIdent()
}

func (m *TypeEmbedded) IsOrdIdent() bool {
	return m.Kind == EtkOrdIdent
}

func (m *TypeEmbedded) IsStringType() bool {
	return m.Kind == EtkString
}

var embeddedRealTypeDeclMap = map[string]*astcore.Decl{}
var embeddedOrdIdentTypeDeclMap = map[string]*astcore.Decl{}
var embeddedStringTypeDeclMap = map[string]*astcore.Decl{}

func newEmbeddedTypeDecl(kind EmbeddedTypeKind, name string) *TypeDecl {
	typ := newTypeEmbedded(kind, name)
	typeDecl := &TypeDecl{Ident: typ.Ident, Type: typ}
	key := strings.ToUpper(typ.Ident.Name)
	decl := typeDecl.ToDeclarations()[0]
	switch kind {
	case EtkReal:
		embeddedRealTypeDeclMap[key] = decl
	case EtkOrdIdent:
		embeddedOrdIdentTypeDeclMap[key] = decl
	case EtkString:
		embeddedStringTypeDeclMap[key] = decl
	default:
		panic(errors.Errorf("unexpected embedded type kind"))
	}
	return typeDecl
}

func EmbeddedRealTypeDecl(name string) *astcore.Decl {
	return embeddedRealTypeDeclMap[strings.ToUpper(name)]
}
func EmbeddedOrdIdentTypeDecl(name string) *astcore.Decl {
	return embeddedOrdIdentTypeDeclMap[strings.ToUpper(name)]
}
func EmbeddedStringTypeDecl(name string) *astcore.Decl {
	return embeddedStringTypeDeclMap[strings.ToUpper(name)]
}

var (
	EmbeddedReal48   = newEmbeddedTypeDecl(EtkReal, "Real48")
	EmbeddedReal     = newEmbeddedTypeDecl(EtkReal, "Real")
	EmbeddedSingile  = newEmbeddedTypeDecl(EtkReal, "Single")
	EmbeddedDouble   = newEmbeddedTypeDecl(EtkReal, "Double")
	EmbeddedExtended = newEmbeddedTypeDecl(EtkReal, "Extended")
	EmbeddedCurrency = newEmbeddedTypeDecl(EtkReal, "Currency")
	EmbeddedComp     = newEmbeddedTypeDecl(EtkReal, "Comp")

	EmbeddedInteger  = newEmbeddedTypeDecl(EtkOrdIdent, "Integer")
	EmbeddedCardinal = newEmbeddedTypeDecl(EtkOrdIdent, "Cardinal")
	EmbeddedShortInt = newEmbeddedTypeDecl(EtkOrdIdent, "ShortInt")
	EmbeddedSmallInt = newEmbeddedTypeDecl(EtkOrdIdent, "SmallInt")
	EmbeddedLongInt  = newEmbeddedTypeDecl(EtkOrdIdent, "LongInt")
	EmbeddedInt64    = newEmbeddedTypeDecl(EtkOrdIdent, "Int64")
	EmbeddedByte     = newEmbeddedTypeDecl(EtkOrdIdent, "Byte")
	EmbeddedWord     = newEmbeddedTypeDecl(EtkOrdIdent, "Word")
	EmbeddedLongWord = newEmbeddedTypeDecl(EtkOrdIdent, "LongWord")
	EmbeddedChar     = newEmbeddedTypeDecl(EtkOrdIdent, "Char")
	EmbeddedAnsiChar = newEmbeddedTypeDecl(EtkOrdIdent, "AnsiChar")
	EmbeddedWideChar = newEmbeddedTypeDecl(EtkOrdIdent, "WideChar")
	EmbeddedBoolean  = newEmbeddedTypeDecl(EtkOrdIdent, "Boolean")

	EmbeddedString     = newEmbeddedTypeDecl(EtkString, "String")
	EmbeddedAnsiString = newEmbeddedTypeDecl(EtkString, "AnsiString")
	EmbeddedWideString = newEmbeddedTypeDecl(EtkString, "WideString")

	// TODO Define Embedded Pointer Types PChar, PInteger, PByteArray, etc.
)
