package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertDirectives(t *testing.T, expected, actual []ast.Directive) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertDirective(t, exp, act)
		}
	}
}

func AssertDirective(t *testing.T, expected, actual ast.Directive) {
	// Do nothing
}

func AssertExternalOptions(t *testing.T, expected, actual *ast.ExternalOptions) {
	assert.Equal(t, expected.LibraryName, actual.LibraryName)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Index, actual.Index)
}
