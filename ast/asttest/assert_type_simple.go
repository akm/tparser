package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertSimpleType(t *testing.T, expected, actual ast.SimpleType) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.RealType:
		AssertRealType(t, exp, actual.(*ast.RealType))
	case ast.OrdinalType:
		AssertOrdinalType(t, exp, actual.(ast.OrdinalType))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertRealType(t *testing.T, expected, actual *ast.RealType) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
}

func AssertOrdinalType(t *testing.T, expected, actual ast.OrdinalType) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.OrdIdent:
		AssertOrdIdent(t, exp, actual.(*ast.OrdIdent))
	case ast.EnumeratedType:
		AssertEnumeratedType(t, exp, actual.(ast.EnumeratedType))
	case *ast.SubrangeType:
		AssertSubrangeType(t, exp, actual.(*ast.SubrangeType))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertOrdIdent(t *testing.T, expected, actual *ast.OrdIdent) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
}

func AssertEnumeratedType(t *testing.T, expected, actual ast.EnumeratedType) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertEnumeratedTypeElement(t, exp, act)
		}
	}
}

func AssertEnumeratedTypeElement(t *testing.T, expected, actual *ast.EnumeratedTypeElement) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.ConstExpr, actual.ConstExpr) {
		AssertConstExpr(t, expected.ConstExpr, actual.ConstExpr)
	}
}

func AssertSubrangeType(t *testing.T, expected, actual *ast.SubrangeType) {
	if !assert.Equal(t, expected.Low, actual.Low) {
		AssertConstExpr(t, expected.Low, actual.Low)
	}
	if !assert.Equal(t, expected.High, actual.High) {
		AssertConstExpr(t, expected.High, actual.High)
	}
}
