package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestStrucArray(t *testing.T) {
	NewTypeTestRunner(t,
		"Char array",
		[]rune(`array[1..100] of Char`),
		&ast.ArrayType{
			IndexTypes: []ast.OrdinalType{
				&ast.SubrangeType{
					Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
					High: asttest.NewConstExpr(asttest.NewNumber("100")),
				},
			},
			BaseType: &ast.OrdIdent{Ident: asttest.NewIdent("Char")},
		},
	).Run().RunVarSection("MyArray")

	NewTypeTestRunner(t,
		"Matrix by array of array",
		[]rune(`array[1..10] of array[1..50] of Real`),
		&ast.ArrayType{
			IndexTypes: []ast.OrdinalType{
				&ast.SubrangeType{
					Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
					High: asttest.NewConstExpr(asttest.NewNumber("10")),
				},
			},
			BaseType: &ast.ArrayType{
				IndexTypes: []ast.OrdinalType{
					&ast.SubrangeType{
						Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
						High: asttest.NewConstExpr(asttest.NewNumber("50")),
					},
				},
				BaseType: &ast.RealType{Ident: asttest.NewIdent("Real")},
			},
		},
	).Run().RunTypeSection("TMatrix")

	NewTypeTestRunner(t,
		"Matrix by array with 2 indexes",
		[]rune(`array[1..10, 1..50] of Real`),
		&ast.ArrayType{
			IndexTypes: []ast.OrdinalType{
				&ast.SubrangeType{
					Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
					High: asttest.NewConstExpr(asttest.NewNumber("10")),
				},
				&ast.SubrangeType{
					Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
					High: asttest.NewConstExpr(asttest.NewNumber("50")),
				},
			},
			BaseType: &ast.RealType{Ident: asttest.NewIdent("Real")},
		},
	).Run().RunTypeSection("TMatrix")

	tshoeSizeDecl := &ast.TypeDecl{
		Ident: asttest.NewIdent("TShoeSize"),
		Type: &ast.SubrangeType{
			Low:  asttest.NewConstExpr(asttest.NewNumber("24")),
			High: asttest.NewConstExpr(asttest.NewNumber("27")),
		},
	}
	tshoeSizeContext := func() interface{} {
		declMap := astcore.NewDeclMap()
		assert.NoError(t, declMap.Set(tshoeSizeDecl))
		r := parser.NewContext(declMap)
		return r
	}

	NewTypeTestRunner(t,
		"array with 3 complicated indexes",
		[]rune(`packed array[Boolean,1..10,TShoeSize] of Integer`),
		&ast.ArrayType{
			Packed: true,
			IndexTypes: []ast.OrdinalType{
				&ast.OrdIdent{Ident: asttest.NewIdent("Boolean")},
				&ast.SubrangeType{
					Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
					High: asttest.NewConstExpr(asttest.NewNumber("10")),
				},
				&ast.TypeId{
					Ident: asttest.NewIdent("TShoeSize"),
					Ref:   tshoeSizeDecl.ToDeclarations()[0],
				},
			},
			BaseType: &ast.OrdIdent{Ident: asttest.NewIdent("Integer")},
		},
		tshoeSizeContext,
	).Run()

	NewTypeTestRunner(t,
		"nested arrays",
		[]rune(`packed array[Boolean] of packed array[1..10] of packed array[TShoeSize] of Integer`),
		&ast.ArrayType{
			Packed: true,
			IndexTypes: []ast.OrdinalType{
				&ast.OrdIdent{Ident: asttest.NewIdent("Boolean")},
			},
			BaseType: &ast.ArrayType{
				Packed: true,
				IndexTypes: []ast.OrdinalType{
					&ast.SubrangeType{
						Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
						High: asttest.NewConstExpr(asttest.NewNumber("10")),
					},
				},
				BaseType: &ast.ArrayType{
					Packed: true,
					IndexTypes: []ast.OrdinalType{
						&ast.TypeId{
							Ident: asttest.NewIdent("TShoeSize"),
							Ref:   tshoeSizeDecl.ToDeclarations()[0],
						},
					},
					BaseType: &ast.OrdIdent{Ident: asttest.NewIdent("Integer")},
				},
			},
		},
		tshoeSizeContext,
	).Run()

	NewTypeTestRunner(t,
		"dynamic arrays",
		[]rune(`array of Real`),
		&ast.ArrayType{
			IndexTypes: nil,
			BaseType:   &ast.RealType{Ident: asttest.NewIdent("Real")},
		},
	).Run().RunVarSection("MyFlexibleArray")

	NewTypeTestRunner(t,
		"multidementional dynamic arrays",
		[]rune(`array of array of string`),
		&ast.ArrayType{
			IndexTypes: nil,
			BaseType: &ast.ArrayType{
				IndexTypes: nil,
				BaseType:   &ast.StringType{Name: "STRING"},
			},
		},
	).Run().RunTypeSection("TMessageGrid")
}
