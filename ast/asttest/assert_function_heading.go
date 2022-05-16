package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertExportedHeading(t *testing.T, expected, actual *ast.ExportedHeading) {
	if !assert.Equal(t, expected.FunctionHeading, actual.FunctionHeading) {
		AssertFunctionHeading(t, expected.FunctionHeading, actual.FunctionHeading)
	}
	if !assert.Equal(t, expected.Directives, actual.Directives) {
		AssertDirectives(t, expected.Directives, actual.Directives)
	}
	if !assert.Equal(t, expected.ExternalOptions, actual.ExternalOptions) {
		AssertExternalOptions(t, expected.ExternalOptions, actual.ExternalOptions)
	}
}

func AssertFunctionHeading(t *testing.T, expected, actual *ast.FunctionHeading) {
	assert.Equal(t, expected.Type, actual.Type)
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.FormalParameters, actual.FormalParameters) {
		AssertFormalParameters(t, expected.FormalParameters, actual.FormalParameters)
	}
	if !assert.Equal(t, expected.ReturnType, actual.ReturnType) {
		AssertType(t, expected.ReturnType, actual.ReturnType)
	}
}

func AssertFormalParameters(t *testing.T, expected, actual ast.FormalParameters) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertFormalParm(t, exp, act)
		}
	}
}

func AssertFormalParm(t *testing.T, expected, actual *ast.FormalParm) {
	assert.Equal(t, expected.Opt, actual.Opt)
	if !assert.Equal(t, expected.Parameter, actual.Parameter) {
		AssertParameter(t, expected.Parameter, actual.Parameter)
	}
}

func AssertParameterType(t *testing.T, expected, actual *ast.ParameterType) {
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
	assert.Equal(t, expected.IsArray, actual.IsArray)
}

func AssertParameter(t *testing.T, expected, actual *ast.Parameter) {
	if !assert.Equal(t, expected.IdentList, actual.IdentList) {
		AssertIdentList(t, expected.IdentList, actual.IdentList)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertParameterType(t, expected.Type, actual.Type)
	}
	if !assert.Equal(t, expected.ConstExpr, actual.ConstExpr) {
		AssertConstExpr(t, expected.ConstExpr, actual.ConstExpr)
	}
}
