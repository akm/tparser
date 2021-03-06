package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithTypeSection(t *testing.T) {
	unit2 := asttest.NewUnitId("Unit2")
	RunUnitTest(t,
		"2 type declarations",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		USES Unit2;
		TYPE
			TTypeId1 = TType1;
			TTypeId2 = Unit2.TType2;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				UsesClause: ast.UsesClause{asttest.NewUnitRef("Unit2")},
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: asttest.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: asttest.NewIdent("TType1")}},
						{Ident: asttest.NewIdent("TTypeId2"), Type: &ast.TypeId{UnitId: unit2, Ident: asttest.NewIdent("TType2")}},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)

	RunUnitTest(t,
		"2 type sections",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		USES Unit2;
		TYPE
			TTypeId1 = TType1;
			TTypeId2 = Unit2.TType2;
		TYPE
			TMyInteger1 = INTEGER;
			TMyReal1 = REAL;
			TMyString1 = STRING;
			TMyString2 = ANSISTRING;
			TMyEnumerated1 = (tsClick, tsClack, tsClock);
			TMySubrange1 = tsClick..tsClack;
			TMySubrange2 = -128..127;
			TMySubrange3 = 'A'..'Z';
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				UsesClause: ast.UsesClause{asttest.NewUnitRef("Unit2")},
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: asttest.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: asttest.NewIdent("TType1")}},
						{Ident: asttest.NewIdent("TTypeId2"), Type: &ast.TypeId{UnitId: unit2, Ident: asttest.NewIdent("TType2")}},
					},
					func() ast.TypeSection {
						tsClick := &ast.EnumeratedTypeElement{Ident: asttest.NewIdent("tsClick")}
						tsClack := &ast.EnumeratedTypeElement{Ident: asttest.NewIdent("tsClack")}
						tsClock := &ast.EnumeratedTypeElement{Ident: asttest.NewIdent("tsClock")}
						return ast.TypeSection{
							{Ident: asttest.NewIdent("TMyInteger1"), Type: ast.NewOrdIdent(asttest.NewIdent("INTEGER"))},
							{Ident: asttest.NewIdent("TMyReal1"), Type: ast.NewRealType(asttest.NewIdent("REAL"))},
							{Ident: asttest.NewIdent("TMyString1"), Type: asttest.NewStringType("STRING")},
							{Ident: asttest.NewIdent("TMyString2"), Type: asttest.NewStringType("ANSISTRING")},
							{Ident: asttest.NewIdent("TMyEnumerated1"), Type: ast.EnumeratedType{tsClick, tsClack, tsClock}},
							{Ident: asttest.NewIdent("TMySubrange1"), Type: &ast.SubrangeType{
								Low:  asttest.NewConstExpr(asttest.NewQualId("tsClick", tsClick.ToDeclarations()[0])),
								High: asttest.NewConstExpr(asttest.NewQualId("tsClack", tsClack.ToDeclarations()[0])),
							}},
							{Ident: asttest.NewIdent("TMySubrange2"), Type: &ast.SubrangeType{Low: asttest.NewConstExpr(asttest.NewNumber("-128")), High: asttest.NewConstExpr(asttest.NewNumber("127"))}},
							{Ident: asttest.NewIdent("TMySubrange3"), Type: &ast.SubrangeType{Low: asttest.NewConstExpr(asttest.NewString("'A'")), High: asttest.NewConstExpr(asttest.NewString("'Z'"))}},
						}
					}(),
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestTypeSection(t *testing.T) {
	RunTypeSection(t,
		"simple type declaration",
		[]rune(`TYPE TTypeId1 = TType1;`),
		ast.TypeSection{
			{Ident: asttest.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: asttest.NewIdent("TType1")}},
		},
	)
	RunTypeSection(t,
		"2 type declarations",
		[]rune(`TYPE TTypeId1 = TType1;
			TTypeId2 = (tsClick, tsClack, tsClock);`),
		ast.TypeSection{
			{Ident: asttest.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: asttest.NewIdent("TType1")}},
			{Ident: asttest.NewIdent("TTypeId2"), Type: ast.EnumeratedType{
				{Ident: asttest.NewIdent("tsClick")},
				{Ident: asttest.NewIdent("tsClack")},
				{Ident: asttest.NewIdent("tsClock")},
			}},
		},
	)
	RunTypeSection(t,
		"type declaration with RealType",
		[]rune(`TYPE TRealType1 = REAL;`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TRealType1"),
				Type:  asttest.NewRealType("REAL"),
			},
		},
	)
	RunTypeSection(t,
		"several type declaration",
		[]rune(`TYPE
			TMyInteger1 = INTEGER;
			TMyString1 = STRING;
			TMyReal1 = REAL;`),
		ast.TypeSection{
			{Ident: asttest.NewIdent("TMyInteger1"), Type: asttest.NewOrdIdent(asttest.NewIdent("INTEGER"))},
			{Ident: asttest.NewIdent("TMyString1"), Type: asttest.NewStringType("STRING")},
			{Ident: asttest.NewIdent("TMyReal1"), Type: asttest.NewRealType(asttest.NewIdent("REAL"))},
		},
	)
}

func newDeclMapWithU1() astcore.DeclMap {
	u1 := &ast.Unit{
		Ident: asttest.NewIdent("U1"),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				ast.TypeSection{
					{Ident: asttest.NewIdent("TType1"), Type: &ast.TypeId{Ident: asttest.NewIdent("String")}},
				},
			},
		},
	}

	usesClauseItemToU1 := &ast.UsesClauseItem{
		Ident: asttest.NewIdent("U1"),
		Unit:  u1,
	}

	declMap := astcore.DeclMapImpl{"u1": usesClauseItemToU1.ToDeclarations()[0]}
	return declMap
}

func TestTypeDecl(t *testing.T) {

	run := func(name string, text []rune, expected *ast.TypeDecl) {
		t.Run(name, func(t *testing.T) {
			declMap := newDeclMapWithU1()
			assert.NotNil(t, declMap.Get("U1"))
			parser := NewTestUnitParser(&text, NewTestUnitContext(declMap))
			parser.NextToken()
			res, err := parser.ParseTypeDecl()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"simple type",
		[]rune(`TTypeId1 = TType1`),
		&ast.TypeDecl{
			Ident: asttest.NewIdent("TTypeId1", asttest.NewIdentLocation(1, 1, 0, 1, 9, 8)),
			Type:  &ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 12, 11, 1, 18, 17))},
		},
	)

	run(
		"simple type with unit",
		[]rune(`TTypeId1 = U1.TType1`),
		&ast.TypeDecl{
			Ident: asttest.NewIdent("TTypeId1", asttest.NewIdentLocation(1, 1, 0, 1, 9, 8)),
			Type: &ast.TypeId{
				UnitId: ast.NewUnitId(asttest.NewIdent("U1", asttest.NewIdentLocation(1, 12, 11, 14))),
				Ident:  asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 15, 14, 1, 21, 20)),
			},
		},
	)
	run(
		"type declaration with TYPE reserved word",
		[]rune(`TTypeId1 = TYPE TType1`),
		&ast.TypeDecl{
			Ident: asttest.NewIdent("TTypeId1", asttest.NewIdentLocation(1, 1, 0, 1, 9, 8)),
			Type:  &ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 17, 16, 1, 23, 22))},
		},
	)
}

func TestTypeId(t *testing.T) {
	declMap := newDeclMapWithU1()
	assert.NotNil(t, declMap.Get("U1"))

	run := func(name string, text []rune, expected *ast.TypeId) {
		t.Run(name, func(t *testing.T) {
			parser := NewTestParser(&text, NewTestUnitContext(declMap))
			parser.NextToken()
			res, err := parser.ParseTypeForIdentifier()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"simple type",
		[]rune(`TType1`),
		&ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 1, 0, 1, 7, 6))},
	)

	run(
		"type with unit",
		[]rune(`U1.TType1`),
		&ast.TypeId{
			UnitId: ast.NewUnitId(asttest.NewIdent("U1", asttest.NewIdentLocation(1, 1, 0, 3))),
			Ident:  asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 4, 3, 1, 10, 9)),
		},
	)
}

func TestNamedType(t *testing.T) {
	run := func(text []rune, expected ast.Type) {
		RunTypeTest(t, string(text), text, expected)
	}

	run([]rune(`INTEGER`), ast.NewOrdIdent(asttest.NewIdent("INTEGER")))
	run([]rune(`CARDINAL`), ast.NewOrdIdent(asttest.NewIdent("CARDINAL")))
	run([]rune(`SHORTINT`), ast.NewOrdIdent(asttest.NewIdent("SHORTINT")))
	run([]rune(`SMALLINT`), ast.NewOrdIdent(asttest.NewIdent("SMALLINT")))
	run([]rune(`LONGINT`), ast.NewOrdIdent(asttest.NewIdent("LONGINT")))
	run([]rune(`INT64`), ast.NewOrdIdent(asttest.NewIdent("INT64")))
	run([]rune(`BYTE`), ast.NewOrdIdent(asttest.NewIdent("BYTE")))
	run([]rune(`WORD`), ast.NewOrdIdent(asttest.NewIdent("WORD")))
	run([]rune(`LONGWORD`), ast.NewOrdIdent(asttest.NewIdent("LONGWORD")))
	run([]rune(`CHAR`), ast.NewOrdIdent(asttest.NewIdent("CHAR")))
	run([]rune(`ANSICHAR`), ast.NewOrdIdent(asttest.NewIdent("ANSICHAR")))
	run([]rune(`WIDECHAR`), ast.NewOrdIdent(asttest.NewIdent("WIDECHAR")))
	run([]rune(`BOOLEAN`), ast.NewOrdIdent(asttest.NewIdent("BOOLEAN")))

	run([]rune(`REAL48`), ast.NewRealType(asttest.NewIdent("REAL48")))
	run([]rune(`REAL`), ast.NewRealType(asttest.NewIdent("REAL")))
	run([]rune(`SINGLE`), ast.NewRealType(asttest.NewIdent("SINGLE")))
	run([]rune(`DOUBLE`), ast.NewRealType(asttest.NewIdent("DOUBLE")))
	run([]rune(`EXTENDED`), ast.NewRealType(asttest.NewIdent("EXTENDED")))
	run([]rune(`CURRENCY`), ast.NewRealType(asttest.NewIdent("CURRENCY")))
	run([]rune(`COMP`), ast.NewRealType(asttest.NewIdent("COMP")))
}
