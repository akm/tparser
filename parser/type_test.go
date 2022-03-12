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

	u2 := ast.UnitId("U2")
	u4 := ast.UnitId("U4")
	run(
		"2 type declarations",
		[]rune(`
		UNIT U1;
		INTERFACE
		TYPE
			TTypeId1 = TType1;
			TTypeId2 = U2.TType2;
		TYPE TTypeId3 = TType3;
		TYPE TTypeId4 = U4.TType4;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: ast.Ident("U1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.TypeSection{
						{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
						{Ident: ast.Ident("TTypeId2"), Type: &ast.TypeId{UnitId: &u2, Ident: ast.Ident("TType2")}},
					},
					ast.TypeSection{
						{Ident: ast.Ident("TTypeId3"), Type: &ast.TypeId{Ident: ast.Ident("TType3")}},
					},
					ast.TypeSection{
						{Ident: ast.Ident("TTypeId4"), Type: &ast.TypeId{UnitId: &u4, Ident: ast.Ident("TType4")}},
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
			TTypeId2 = TType2;`),
		ast.TypeSection{
			{Ident: ast.Ident("TTypeId1"), Type: &ast.TypeId{Ident: ast.Ident("TType1")}},
			{Ident: ast.Ident("TTypeId2"), Type: &ast.TypeId{Ident: ast.Ident("TType2")}},
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
			res, err := parser.ParseTypeId()
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
