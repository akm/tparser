package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertStringType(t *testing.T, expected, actual ast.StringType) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.TypeEmbedded:
		AssertTypeEmbedded(t, exp, actual.(*ast.TypeEmbedded))
	case *ast.TypeId:
		AssertTypeId(t, exp, actual.(*ast.TypeId))
	case *ast.FixedStringType:
		AssertFixedStringType(t, exp, actual.(*ast.FixedStringType))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertFixedStringType(t *testing.T, expected, actual *ast.FixedStringType) {
	if !assert.Equal(t, expected.StringType, actual.StringType) {
		AssertStringType(t, expected.StringType, actual.StringType)
	}
	AssertConstExpr(t, expected.Length, actual.Length)
}
