package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertStrucType(t *testing.T, expected, actual ast.StrucType) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.ArrayType:
		AssertArrayType(t, exp, actual.(*ast.ArrayType))
	case *ast.SetType:
		AssertSetType(t, exp, actual.(*ast.SetType))
	case *ast.FileType:
		AssertFileType(t, exp, actual.(*ast.FileType))
	case *ast.RecType:
		AssertRecType(t, exp, actual.(*ast.RecType))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertArrayType(t *testing.T, expected, actual *ast.ArrayType) {
	if !assert.Equal(t, expected.IndexTypes, actual.IndexTypes) {
		AssertOrdinalTypes(t, expected.IndexTypes, actual.IndexTypes)
	}
	if !assert.Equal(t, expected.BaseType, actual.BaseType) {
		AssertType(t, expected.BaseType, actual.BaseType)
	}
	assert.Equal(t, expected.Packed, actual.Packed)
}

func AssertSetType(t *testing.T, expected, actual *ast.SetType) {
	if !assert.Equal(t, expected.OrdinalType, actual.OrdinalType) {
		AssertOrdinalType(t, expected.OrdinalType, actual.OrdinalType)
	}
	assert.Equal(t, expected.Packed, actual.Packed)
}

func AssertFileType(t *testing.T, expected, actual *ast.FileType) {
	if !assert.Equal(t, expected.TypeId, actual.TypeId) {
		AssertTypeId(t, expected.TypeId, actual.TypeId)
	}
	assert.Equal(t, expected.Packed, actual.Packed)
}

func AssertRecType(t *testing.T, expected, actual *ast.RecType) {
	if !assert.Equal(t, expected.FieldList, actual.FieldList) {
		AssertFieldList(t, expected.FieldList, actual.FieldList)
	}
	assert.Equal(t, expected.Packed, actual.Packed)
}

func AssertFieldList(t *testing.T, expected, actual *ast.FieldList) {
	if !assert.Equal(t, expected.FieldDecls, actual.FieldDecls) {
		AssertFieldDecls(t, expected.FieldDecls, actual.FieldDecls)
	}
	if !assert.Equal(t, expected.VariantSection, actual.VariantSection) {
		AssertVariantSection(t, expected.VariantSection, actual.VariantSection)
	}
}

func AssertFieldDecls(t *testing.T, expected, actual ast.FieldDecls) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertFieldDecl(t, exp, act)
		}
	}
}

func AssertFieldDecl(t *testing.T, expected, actual *ast.FieldDecl) {
	if !assert.Equal(t, expected.IdentList, actual.IdentList) {
		AssertIdentList(t, expected.IdentList, actual.IdentList)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
}

func AssertVariantSection(t *testing.T, expected, actual *ast.VariantSection) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.TypeId, actual.TypeId) {
		AssertOrdinalType(t, expected.TypeId, actual.TypeId)
	}
	if !assert.Equal(t, expected.RecVariants, actual.RecVariants) {
		AssertRecVariants(t, expected.RecVariants, actual.RecVariants)
	}
}

func AssertRecVariants(t *testing.T, expected, actual ast.RecVariants) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertRecVariant(t, exp, act)
		}
	}
}

func AssertRecVariant(t *testing.T, expected, actual *ast.RecVariant) {
	if !assert.Equal(t, expected.ConstExprs, actual.ConstExprs) {
		AssertConstExprs(t, expected.ConstExprs, actual.ConstExprs)
	}
	if !assert.Equal(t, expected.FieldList, actual.FieldList) {
		AssertFieldList(t, expected.FieldList, actual.FieldList)
	}
}
