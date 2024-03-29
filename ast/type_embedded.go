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
	EtkStringType
	EtkPointerType
	EtkVariantType
)

var embeddedTypeKindAll = []EmbeddedTypeKind{
	EtkReal,
	EtkOrdIdent,
	EtkStringType,
	EtkPointerType,
	EtkVariantType,
}

type TypeEmbedded struct {
	Kind  EmbeddedTypeKind
	Ident *Ident
}

var _ Type = (*TypeEmbedded)(nil)
var _ SimpleType = (*TypeEmbedded)(nil)
var _ RealType = (*TypeEmbedded)(nil)
var _ OrdinalType = (*TypeEmbedded)(nil)
var _ OrdIdent = (*TypeEmbedded)(nil)
var _ StringType = (*TypeEmbedded)(nil)
var _ PointerType = (*TypeEmbedded)(nil)
var _ VariantType = (*TypeEmbedded)(nil)

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
	return m.Kind == EtkStringType
}

func (m *TypeEmbedded) IsPointerType() bool {
	return m.Kind == EtkPointerType
}

func (m *TypeEmbedded) IsVariantType() bool {
	return m.Kind == EtkVariantType
}

var embeddedTypeDeclMaps = func() map[EmbeddedTypeKind]map[string]*astcore.Decl {
	r := make(map[EmbeddedTypeKind]map[string]*astcore.Decl)
	for _, kind := range embeddedTypeKindAll {
		r[kind] = make(map[string]*astcore.Decl)
	}
	return r
}()

var embeddedTypeDeclMap = map[string]*astcore.Decl{}

func newEmbeddedTypeDecl(kind EmbeddedTypeKind, name string) *TypeDecl {
	typ := newTypeEmbedded(kind, name)
	typeDecl := &TypeDecl{Ident: typ.Ident, Type: typ}
	key := strings.ToUpper(typ.Ident.Name)
	decl := typeDecl.ToDeclarations()[0]
	embeddedTypeDeclMap[key] = decl
	embeddedTypeDeclMaps[kind][key] = decl
	return typeDecl
}

func EmbeddedTypeDecl(kind EmbeddedTypeKind, name string) *astcore.Decl {
	return embeddedTypeDeclMaps[kind][strings.ToUpper(name)]
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

	EmbeddedString     = newEmbeddedTypeDecl(EtkStringType, "String")
	EmbeddedAnsiString = newEmbeddedTypeDecl(EtkStringType, "AnsiString")
	EmbeddedWideString = newEmbeddedTypeDecl(EtkStringType, "WideString")

	EmbeddedPointer   = newEmbeddedTypeDecl(EtkPointerType, "Pointer")
	EmbeddedPChar     = newEmbeddedTypeDecl(EtkPointerType, "PChar")
	EmbeddedPAnsiChar = newEmbeddedTypeDecl(EtkPointerType, "PAnsiChar")
	EmbeddedPWideChar = newEmbeddedTypeDecl(EtkPointerType, "PWideChar")

	EmbeddedVariant    = newEmbeddedTypeDecl(EtkVariantType, "Variant")
	EmbeddedOleVariant = newEmbeddedTypeDecl(EtkVariantType, "OleVariant")
)

type embeddedTypeDeclMapSingleton struct {
}

func (m *embeddedTypeDeclMapSingleton) Get(name string) *astcore.Decl {
	return embeddedTypeDeclMap[strings.ToUpper(name)]
}

func (m *embeddedTypeDeclMapSingleton) Set(astcore.DeclNode) error {
	return errors.Errorf("Can't set anything to embedded type decl map")
}

func (m *embeddedTypeDeclMapSingleton) Overwrite(name string, decl *astcore.Decl) {
}

var EmbeddedTypeDeclMap = &embeddedTypeDeclMapSingleton{}
