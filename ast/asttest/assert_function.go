package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertProcedureDeclSection(t *testing.T, expected, actual ast.ProcedureDeclSection) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.FunctionDecl:
		AssertFunctionDecl(t, exp, actual.(*ast.FunctionDecl))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertFunctionDecl(t *testing.T, expected, actual *ast.FunctionDecl) {
	if !assert.Equal(t, expected.FunctionHeading, actual.FunctionHeading) {
		AssertFunctionHeading(t, expected.FunctionHeading, actual.FunctionHeading)
	}
	if !assert.Equal(t, expected.Directives, actual.Directives) {
		AssertDirectives(t, expected.Directives, actual.Directives)
	}
	if !assert.Equal(t, expected.ExternalOptions, actual.ExternalOptions) {
		AssertExternalOptions(t, expected.ExternalOptions, actual.ExternalOptions)
	}
	if !assert.Equal(t, expected.PortabilityDirective, actual.PortabilityDirective) {
		AssertPortabilityDirective(t, expected.PortabilityDirective, actual.PortabilityDirective)
	}
	if !assert.Equal(t, expected.Block, actual.Block) {
		AssertBlock(t, expected.Block, actual.Block)
	}
}
