package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
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
			Ident: ast.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				UsesClause: &ast.UsesClause{ast.NewUnitRef("Unit2")},
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: ast.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: ast.NewIdent("TType1")}},
						{Ident: ast.NewIdent("TTypeId2"), Type: &ast.TypeId{UnitId: &unit2, Ident: ast.NewIdent("TType2")}},
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
			Ident: ast.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				UsesClause: &ast.UsesClause{ast.NewUnitRef("Unit2")},
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: ast.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: ast.NewIdent("TType1")}},
						{Ident: ast.NewIdent("TTypeId2"), Type: &ast.TypeId{UnitId: &unit2, Ident: ast.NewIdent("TType2")}},
					},
					ast.TypeSection{
						{Ident: ast.NewIdent("TMyInteger1"), Type: &ast.OrdIdent{Name: ast.NewIdent("INTEGER")}},
						{Ident: ast.NewIdent("TMyReal1"), Type: &ast.RealType{Name: ast.NewIdent("REAL")}},
						{Ident: ast.NewIdent("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
						{Ident: ast.NewIdent("TMyString2"), Type: &ast.StringType{Name: "ANSISTRING"}},
						{Ident: ast.NewIdent("TMyEnumerated1"), Type: ast.EnumeratedType{
							{Ident: ast.NewIdent("tsClick")},
							{Ident: ast.NewIdent("tsClack")},
							{Ident: ast.NewIdent("tsClock")},
						}},
						{Ident: ast.NewIdent("TMySubrange1"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr("tsClick"), High: *ast.NewConstExpr("tsClack")}},
						{Ident: ast.NewIdent("TMySubrange2"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewNumber("-128")), High: *ast.NewConstExpr(ast.NewNumber("127"))}},
						{Ident: ast.NewIdent("TMySubrange3"), Type: &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewString("'A'")), High: *ast.NewConstExpr(ast.NewString("'Z'"))}},
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
			{Ident: ast.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: ast.NewIdent("TType1")}},
		},
	)
	run(
		"2 type declarations",
		[]rune(`TYPE TTypeId1 = TType1;
			TTypeId2 = (tsClick, tsClack, tsClock);`),
		ast.TypeSection{
			{Ident: ast.NewIdent("TTypeId1"), Type: &ast.TypeId{Ident: ast.NewIdent("TType1")}},
			{Ident: ast.NewIdent("TTypeId2"), Type: ast.EnumeratedType{
				{Ident: ast.NewIdent("tsClick")},
				{Ident: ast.NewIdent("tsClack")},
				{Ident: ast.NewIdent("tsClock")},
			}},
		},
	)
	run(
		"type declaration with RealType",
		[]rune(`TYPE TRealType1 = REAL;`),
		ast.TypeSection{
			&ast.TypeDecl{
				Ident: ast.NewIdent("TRealType1"),
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
			{Ident: ast.NewIdent("TMyInteger1"), Type: &ast.OrdIdent{Name: ast.NewIdent("INTEGER")}},
			{Ident: ast.NewIdent("TMyString1"), Type: &ast.StringType{Name: "STRING"}},
			{Ident: ast.NewIdent("TMyReal1"), Type: &ast.RealType{Name: ast.NewIdent("REAL")}},
		},
	)
}

func TestTypeDecl(t *testing.T) {
	u1 := ast.UnitId("U1")

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
			Ident: ast.NewIdent("TTypeId1"),
			Type:  &ast.TypeId{Ident: ast.NewIdent("TType1")},
		},
	)

	run(
		"simple type with unit",
		[]rune(`TTypeId1 = U1.TType1`),
		&ast.TypeDecl{
			Ident: ast.NewIdent("TTypeId1"),
			Type:  &ast.TypeId{UnitId: &u1, Ident: ast.NewIdent("TType1")},
		},
	)
	run(
		"type declaration with TYPE reserved word",
		[]rune(`TTypeId1 = TYPE TType1`),
		&ast.TypeDecl{
			Ident: ast.NewIdent("TTypeId1"),
			Type:  &ast.TypeId{Ident: ast.NewIdent("TType1")},
		},
	)
}

func TestTypeId(t *testing.T) {
	u1 := ast.UnitId("U1")

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
		&ast.TypeId{Ident: ast.NewIdent("TType1")},
	)

	run(
		"type with unit",
		[]rune(`U1.TType1`),
		&ast.TypeId{
			UnitId: &u1,
			Ident:  ast.NewIdent("TType1"),
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

	run([]rune(`INTEGER`), &ast.OrdIdent{Name: ast.NewIdent("INTEGER")})
	run([]rune(`CARDINAL`), &ast.OrdIdent{Name: ast.NewIdent("CARDINAL")})
	run([]rune(`SHORTINT`), &ast.OrdIdent{Name: ast.NewIdent("SHORTINT")})
	run([]rune(`SMALLINT`), &ast.OrdIdent{Name: ast.NewIdent("SMALLINT")})
	run([]rune(`LONGINT`), &ast.OrdIdent{Name: ast.NewIdent("LONGINT")})
	run([]rune(`INT64`), &ast.OrdIdent{Name: ast.NewIdent("INT64")})
	run([]rune(`BYTE`), &ast.OrdIdent{Name: ast.NewIdent("BYTE")})
	run([]rune(`WORD`), &ast.OrdIdent{Name: ast.NewIdent("WORD")})
	run([]rune(`LONGWORD`), &ast.OrdIdent{Name: ast.NewIdent("LONGWORD")})
	run([]rune(`CHAR`), &ast.OrdIdent{Name: ast.NewIdent("CHAR")})
	run([]rune(`ANSICHAR`), &ast.OrdIdent{Name: ast.NewIdent("ANSICHAR")})
	run([]rune(`WIDECHAR`), &ast.OrdIdent{Name: ast.NewIdent("WIDECHAR")})
	run([]rune(`BOOLEAN`), &ast.OrdIdent{Name: ast.NewIdent("BOOLEAN")})

	run([]rune(`REAL48`), &ast.RealType{Name: ast.NewIdent("REAL48")})
	run([]rune(`REAL`), &ast.RealType{Name: ast.NewIdent("REAL")})
	run([]rune(`SINGLE`), &ast.RealType{Name: ast.NewIdent("SINGLE")})
	run([]rune(`DOUBLE`), &ast.RealType{Name: ast.NewIdent("DOUBLE")})
	run([]rune(`EXTENDED`), &ast.RealType{Name: ast.NewIdent("EXTENDED")})
	run([]rune(`CURRENCY`), &ast.RealType{Name: ast.NewIdent("CURRENCY")})
	run([]rune(`COMP`), &ast.RealType{Name: ast.NewIdent("COMP")})
}
