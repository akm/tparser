package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertTypeSection(t *testing.T, expected, actual ast.TypeSection) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertTypeDecl(t, exp, act)
		}
	}
}

func AssertTypeDecl(t *testing.T, expected, actual *ast.TypeDecl) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
	if !assert.Equal(t, expected.PortabilityDirective, actual.PortabilityDirective) {
		AssertPortabilityDirective(t, expected.PortabilityDirective, actual.PortabilityDirective)
	}
}

func AssertType(t *testing.T, expected, actual ast.Type) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.TypeId:
		AssertTypeId(t, exp, actual.(*ast.TypeId))
	case ast.SimpleType:
		AssertSimpleType(t, exp, actual.(ast.SimpleType))
	case ast.StrucType:
		AssertStrucType(t, exp, actual.(ast.StrucType))
	// case *ast.PointerType:
	// 	AssertPointerType(t, exp, actual.(*ast.PointerType))
	case *ast.StringType:
		AssertStringType(t, exp, actual.(*ast.StringType))
	// case *ast.ProcedureType:
	// 	AssertProcedureType(t, exp, actual.(*ast.ProcedureType))
	// case *ast.VariantType:
	// 	AssertVariantType(t, exp, actual.(*ast.VariantType))
	// case *ast.ClassRefType:
	// 	AssertClassRefType(t, exp, actual.(*ast.ClassRefType))
	default:
		assert.Fail(t, "unexpected type: %T", exp)
	}
}

func AssertTypeId(t *testing.T, expected, actual *ast.TypeId) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.UnitId, actual.UnitId) {
		AssertUnitId(t, expected.UnitId, actual.UnitId)
	}
	if !assert.Equal(t, expected.Ref, actual.Ref) {
		AssertDeclaration(t, expected.Ref, actual.Ref)
	}
}
