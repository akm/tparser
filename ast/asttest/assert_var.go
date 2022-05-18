package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertVarSection(t *testing.T, expected, actual ast.VarSection) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertVarDecl(t, exp, act)
		}
	}
}

func AssertVarDecl(t *testing.T, expected, actual *ast.VarDecl) {
	if !assert.Equal(t, expected.IdentList, actual.IdentList) {
		AssertIdentList(t, expected.IdentList, actual.IdentList)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
	if !assert.Equal(t, expected.Absolute, actual.Absolute) {
		AssertVarDeclAbsolute(t, expected.Absolute, actual.Absolute)
	}
	if !assert.Equal(t, expected.ConstExpr, actual.ConstExpr) {
		AssertConstExpr(t, expected.ConstExpr, actual.ConstExpr)
	}
	if !assert.Equal(t, expected.PortabilityDirective, actual.PortabilityDirective) {
		AssertPortabilityDirective(t, expected.PortabilityDirective, actual.PortabilityDirective)
	}
}

func AssertVarDeclAbsolute(t *testing.T, expected, actual ast.VarDeclAbsolute) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.VarDeclAbsoluteIdent:
		AssertVarDeclAbsoluteIdent(t, exp, actual.(*ast.VarDeclAbsoluteIdent))
	case *ast.VarDeclAbsoluteConstExpr:
		AssertVarDeclAbsoluteConstExpr(t, exp, actual.(*ast.VarDeclAbsoluteConstExpr))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertVarDeclAbsoluteIdent(t *testing.T, expected, actual *ast.VarDeclAbsoluteIdent) {
	AssertIdent(t, (*ast.Ident)(expected), (*ast.Ident)(actual))
}

func AssertVarDeclAbsoluteConstExpr(t *testing.T, expected, actual *ast.VarDeclAbsoluteConstExpr) {
	AssertConstExpr(t, (*ast.ConstExpr)(expected), (*ast.ConstExpr)(actual))
}

func AssertThreadVarSection(t *testing.T, expected, actual ast.ThreadVarSection) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertThreadVarDecl(t, exp, act)
		}
	}
}

func AssertThreadVarDecl(t *testing.T, expected, actual *ast.ThreadVarDecl) {
	if !assert.Equal(t, expected.IdentList, actual.IdentList) {
		AssertIdentList(t, expected.IdentList, actual.IdentList)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
}
