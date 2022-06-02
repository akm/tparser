package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestStrucSet(t *testing.T) {
	tsomeInts := &ast.TypeDecl{
		Ident: asttest.NewIdent("TSomeInts"),
		Type: &ast.SubrangeType{
			Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
			High: asttest.NewConstExpr(asttest.NewNumber("250")),
		},
	}
	tsomeIntsContext := func() interface{} {
		declMap := astcore.NewDeclMap()
		assert.NoError(t, declMap.Set(tsomeInts))
		r := parser.NewContext(declMap)
		return r
	}

	runType(t,
		"Set with TypeId",
		[]rune(`set of TSomeInts`),
		&ast.SetType{
			OrdinalType: &ast.TypeId{
				Ident: asttest.NewIdent("TSomeInts"),
				Ref:   tsomeInts.ToDeclarations()[0],
			},
		},
		tsomeIntsContext,
	)

	runType(t,
		"Set with subrange type",
		[]rune(`set of 1..250`),
		&ast.SetType{
			OrdinalType: &ast.SubrangeType{
				Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
				High: asttest.NewConstExpr(asttest.NewNumber("250")),
			},
		},
	)

	runType(t,
		"Set with a to z",
		[]rune(`set of 'a'..'z'`),
		&ast.SetType{
			OrdinalType: &ast.SubrangeType{
				Low:  asttest.NewConstExpr(asttest.NewString("'a'")),
				High: asttest.NewConstExpr(asttest.NewString("'z'")),
			},
		},
	)

	runType(t,
		"Set of Byte",
		[]rune(`set of Byte`),
		&ast.SetType{
			OrdinalType: &ast.OrdIdent{Ident: asttest.NewIdent("Byte")},
		},
	)

	runType(t,
		"Set of enumerated type",
		[]rune(`set of (Club, Diamond, Heart, Spade)`),
		&ast.SetType{
			OrdinalType: ast.EnumeratedType{
				{Ident: asttest.NewIdent("Club")},
				{Ident: asttest.NewIdent("Diamond")},
				{Ident: asttest.NewIdent("Heart")},
				{Ident: asttest.NewIdent("Spade")},
			},
		},
	)

	runType(t,
		"Set of Char",
		[]rune(`set of Char`),
		&ast.SetType{
			OrdinalType: &ast.OrdIdent{Ident: asttest.NewIdent("Char")},
		},
	)
}
