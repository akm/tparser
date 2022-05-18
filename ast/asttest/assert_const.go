package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertConstSection(t *testing.T, expected, actual ast.ConstSection) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertConstantDecl(t, exp, act)
		}
	}
}

func AssertConstantDecl(t *testing.T, expected, actual *ast.ConstantDecl) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
	if !assert.Equal(t, expected.ConstExpr, actual.ConstExpr) {
		AssertConstExpr(t, expected.ConstExpr, actual.ConstExpr)
	}
	if !assert.Equal(t, expected.PortabilityDirective, actual.PortabilityDirective) {
		AssertPortabilityDirective(t, expected.PortabilityDirective, actual.PortabilityDirective)
	}
}

func AssertConstExpr(t *testing.T, expected, actual *ast.ConstExpr) {
	AssertExpression(t, expected, actual)
}
