package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithTypeSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	unit2 := ast.UnitId("Unit2")
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
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				UsesClause: &ast.UsesClause{"Unit2"},
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
						{Ident: ast.Ident("TTypeId2"), Type: &ast.TypeId{UnitId: &unit2, Ident: ast.Ident("TType2")}},
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
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
						{Ident: ast.Ident("TTypeId2"), Type: &ast.TypeId{UnitId: &unit2, Ident: ast.Ident("TType2")}},
					},
					ast.TypeSection{
						{Ident: ast.Ident("TMyInteger1"), Type: &ast.OrdIdent{Name: ast.Ident("INTEGER")}},
						{Ident: ast.Ident("TMyReal1"), Type: &ast.RealType{Name: ast.Ident("REAL")}},
						{Ident: ast.Ident("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
						{Ident: ast.Ident("TMyString2"), Type: &ast.StringType{Name: "ANSISTRING"}},
						{Ident: ast.Ident("TMyEnumerated1"), Type: ast.EnumeratedType{
							{Ident: ast.Ident("tsClick")},
							{Ident: ast.Ident("tsClack")},
							{Ident: ast.Ident("tsClock")},
						}},
						{Ident: ast.Ident("TMySubrange1"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr("tsClick"), High: *ast.NewConstExpr("tsClack")}},
						{Ident: ast.Ident("TMySubrange2"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewNumber("-128")), High: *ast.NewConstExpr(ast.NewNumber("127"))}},
						{Ident: ast.Ident("TMySubrange3"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewString("'A'")), High: *ast.NewConstExpr(ast.NewString("'Z'"))}},
					},
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
			res, err := parser.ParseTypeSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run("simple type declaration",
		[]rune(`TYPE TTypeId1 = TType1;`),
		ast.TypeSection{
			{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
		},
	)
	run(
		"2 type declarations",
		[]rune(`TYPE TTypeId1 = TType1;
			TTypeId2 = (tsClick, tsClack, tsClock);`),
		ast.TypeSection{
			{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
			{Ident: ast.Ident("TTypeId2"), Type: ast.EnumeratedType{
				{Ident: ast.Ident("tsClick")},
				{Ident: ast.Ident("tsClack")},
				{Ident: ast.Ident("tsClock")},
			}},
		},
	)
	run(
		"type declaration with RealType",
		[]rune(`TYPE TRealType1 = REAL;`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: ast.Ident("TRealType1"),
				Type:  &ast.RealType{Name: "REAL"},
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
			{Ident: ast.Ident("TMyInteger1"), Type: &ast.OrdIdent{Name: ast.Ident("INTEGER")}},
			{Ident: ast.Ident("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
			{Ident: ast.Ident("TMyReal1"), Type: &ast.RealType{Name: ast.Ident("REAL")}},
		},
	)
}

func TestTypeDecl(t *testing.T) {
	run := func(name string, text []rune, expected *ast.TypeDecl) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
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
			Ident: ast.Ident("TTypeId1"),
			Type:  &ast.TypeId{Ident: ast.Ident("TType1")},
		},
	)
	u1 := ast.UnitId("U1")
	run(
		"simple type with unit",
		[]rune(`TTypeId1 = U1.TType1`),
		&ast.TypeDecl{
			Ident: ast.Ident("TTypeId1"),
			Type:  &ast.TypeId{UnitId: &u1, Ident: ast.Ident("TType1")},
		},
	)
	run(
		"type declaration with TYPE reserved word",
		[]rune(`TTypeId1 = TYPE TType1`),
		&ast.TypeDecl{
			Ident: ast.Ident("TTypeId1"),
			Type:  &ast.TypeId{Ident: ast.Ident("TType1")},
		},
	)
}

func TestTypeId(t *testing.T) {
	run := func(name string, text []rune, expected *ast.TypeId) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseTypeIdOrSubrangeType()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"simple type",
		[]rune(`TType1`),
		&ast.TypeId{Ident: ast.Ident("TType1")},
	)

	u1 := ast.UnitId("U1")
	run(
		"type with unit",
		[]rune(`U1.TType1`),
		&ast.TypeId{
			UnitId: &u1,
			Ident:  ast.Ident("TType1"),
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
				assert.Equal(t, expected, res)
			}
		})
	}

	run([]rune(`INTEGER`), &ast.OrdIdent{Name: ast.Ident("INTEGER")})
	run([]rune(`CARDINAL`), &ast.OrdIdent{Name: ast.Ident("CARDINAL")})
	run([]rune(`SHORTINT`), &ast.OrdIdent{Name: ast.Ident("SHORTINT")})
	run([]rune(`SMALLINT`), &ast.OrdIdent{Name: ast.Ident("SMALLINT")})
	run([]rune(`LONGINT`), &ast.OrdIdent{Name: ast.Ident("LONGINT")})
	run([]rune(`INT64`), &ast.OrdIdent{Name: ast.Ident("INT64")})
	run([]rune(`BYTE`), &ast.OrdIdent{Name: ast.Ident("BYTE")})
	run([]rune(`WORD`), &ast.OrdIdent{Name: ast.Ident("WORD")})
	run([]rune(`LONGWORD`), &ast.OrdIdent{Name: ast.Ident("LONGWORD")})
	run([]rune(`CHAR`), &ast.OrdIdent{Name: ast.Ident("CHAR")})
	run([]rune(`ANSICHAR`), &ast.OrdIdent{Name: ast.Ident("ANSICHAR")})
	run([]rune(`WIDECHAR`), &ast.OrdIdent{Name: ast.Ident("WIDECHAR")})
	run([]rune(`BOOLEAN`), &ast.OrdIdent{Name: ast.Ident("BOOLEAN")})

	run([]rune(`REAL48`), &ast.RealType{Name: ast.Ident("REAL48")})
	run([]rune(`REAL`), &ast.RealType{Name: ast.Ident("REAL")})
	run([]rune(`SINGLE`), &ast.RealType{Name: ast.Ident("SINGLE")})
	run([]rune(`DOUBLE`), &ast.RealType{Name: ast.Ident("DOUBLE")})
	run([]rune(`EXTENDED`), &ast.RealType{Name: ast.Ident("EXTENDED")})
	run([]rune(`CURRENCY`), &ast.RealType{Name: ast.Ident("CURRENCY")})
	run([]rune(`COMP`), &ast.RealType{Name: ast.Ident("COMP")})
}
