package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithTypeSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				if !assert.Equal(t, expected, res) {
					if !assert.Equal(t, expected.InterfaceSection.InterfaceDecls, res.InterfaceSection.InterfaceDecls) {
						assert.Equal(t, expected.InterfaceSection.InterfaceDecls[0], res.InterfaceSection.InterfaceDecls[0])
						assert.Equal(t, expected.InterfaceSection.InterfaceDecls[1], res.InterfaceSection.InterfaceDecls[1])
					}
				}
			}
		})
	}

	unit2 := asttest.NewUnitId("Unit2")
	run(
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
	run(
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
							{Ident: asttest.NewIdent("TMyInteger1"), Type: &ast.OrdIdent{Ident: asttest.NewIdent("INTEGER")}},
							{Ident: asttest.NewIdent("TMyReal1"), Type: &ast.RealType{Ident: asttest.NewIdent("REAL")}},
							{Ident: asttest.NewIdent("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
							{Ident: asttest.NewIdent("TMyString2"), Type: &ast.StringType{Name: "ANSISTRING"}},
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
	run := func(name string, text []rune, expected ast.TypeSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseTypeSection(true)
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run("simple type declaration",
		[]rune(`TYPE TTypeId1 = TType1;`),
		ast.TypeSection{
			{Ident: asttest.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: asttest.NewIdent("TType1")}},
		},
	)
	run(
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
	run(
		"type declaration with RealType",
		[]rune(`TYPE TRealType1 = REAL;`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: asttest.NewIdent("TRealType1"),
				Type:  asttest.NewRealType("REAL"),
			},
		},
	)
	run(
		"several type declaration",
		[]rune(`TYPE
			TMyInteger1 = INTEGER;
			TMyString1 = STRING;
			TMyReal1 = REAL;`),
		ast.TypeSection{
			{Ident: asttest.NewIdent("TMyInteger1"), Type: &ast.OrdIdent{Ident: asttest.NewIdent("INTEGER")}},
			{Ident: asttest.NewIdent("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
			{Ident: asttest.NewIdent("TMyReal1"), Type: &ast.RealType{Ident: asttest.NewIdent("REAL")}},
		},
	)
}

func TestTypeDecl(t *testing.T) {
	u1 := asttest.NewUnitId("U1")

	run := func(name string, text []rune, expected *ast.TypeDecl) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text, NewContext(ext.Strings{u1.String()}))
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
			Type:  &ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 12, 11, 1, 17, 17))},
		},
	)

	run(
		"simple type with unit",
		[]rune(`TTypeId1 = U1.TType1`),
		&ast.TypeDecl{
			Ident: asttest.NewIdent("TTypeId1", asttest.NewIdentLocation(1, 1, 0, 1, 9, 8)),
			Type: &ast.TypeId{
				UnitId: ast.NewUnitId(asttest.NewIdent("U1", asttest.NewIdentLocation(1, 12, 11, 14))),
				Ident:  asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 15, 14, 1, 20, 20)),
			},
		},
	)
	run(
		"type declaration with TYPE reserved word",
		[]rune(`TTypeId1 = TYPE TType1`),
		&ast.TypeDecl{
			Ident: asttest.NewIdent("TTypeId1", asttest.NewIdentLocation(1, 1, 0, 1, 9, 8)),
			Type:  &ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 17, 16, 1, 22, 22))},
		},
	)
}

func TestTypeId(t *testing.T) {
	u1 := asttest.NewUnitId("U1")

	run := func(name string, text []rune, expected *ast.TypeId) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text, NewContext(ext.Strings{u1.String()}))
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
		&ast.TypeId{Ident: asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 1, 0, 1, 6, 6))},
	)

	run(
		"type with unit",
		[]rune(`U1.TType1`),
		&ast.TypeId{
			UnitId: ast.NewUnitId(asttest.NewIdent("U1", asttest.NewIdentLocation(1, 1, 0, 3))),
			Ident:  asttest.NewIdent("TType1", asttest.NewIdentLocation(1, 4, 3, 1, 9, 9)),
		},
	)
}

func TestNamedType(t *testing.T) {
	run := func(text []rune, expected ast.Type) {
		t.Run(string(text), func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseType()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run([]rune(`INTEGER`), &ast.OrdIdent{Ident: asttest.NewIdent("INTEGER")})
	run([]rune(`CARDINAL`), &ast.OrdIdent{Ident: asttest.NewIdent("CARDINAL")})
	run([]rune(`SHORTINT`), &ast.OrdIdent{Ident: asttest.NewIdent("SHORTINT")})
	run([]rune(`SMALLINT`), &ast.OrdIdent{Ident: asttest.NewIdent("SMALLINT")})
	run([]rune(`LONGINT`), &ast.OrdIdent{Ident: asttest.NewIdent("LONGINT")})
	run([]rune(`INT64`), &ast.OrdIdent{Ident: asttest.NewIdent("INT64")})
	run([]rune(`BYTE`), &ast.OrdIdent{Ident: asttest.NewIdent("BYTE")})
	run([]rune(`WORD`), &ast.OrdIdent{Ident: asttest.NewIdent("WORD")})
	run([]rune(`LONGWORD`), &ast.OrdIdent{Ident: asttest.NewIdent("LONGWORD")})
	run([]rune(`CHAR`), &ast.OrdIdent{Ident: asttest.NewIdent("CHAR")})
	run([]rune(`ANSICHAR`), &ast.OrdIdent{Ident: asttest.NewIdent("ANSICHAR")})
	run([]rune(`WIDECHAR`), &ast.OrdIdent{Ident: asttest.NewIdent("WIDECHAR")})
	run([]rune(`BOOLEAN`), &ast.OrdIdent{Ident: asttest.NewIdent("BOOLEAN")})

	run([]rune(`REAL48`), &ast.RealType{Ident: asttest.NewIdent("REAL48")})
	run([]rune(`REAL`), &ast.RealType{Ident: asttest.NewIdent("REAL")})
	run([]rune(`SINGLE`), &ast.RealType{Ident: asttest.NewIdent("SINGLE")})
	run([]rune(`DOUBLE`), &ast.RealType{Ident: asttest.NewIdent("DOUBLE")})
	run([]rune(`EXTENDED`), &ast.RealType{Ident: asttest.NewIdent("EXTENDED")})
	run([]rune(`CURRENCY`), &ast.RealType{Ident: asttest.NewIdent("CURRENCY")})
	run([]rune(`COMP`), &ast.RealType{Ident: asttest.NewIdent("COMP")})
}
