package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertUsesClause(t *testing.T, expected ast.UsesClause, actual ast.UsesClause) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, expRef := range expected {
		actRef := actual[i]
		if !assert.Equal(t, expRef, actRef) {
			AssertUnitRef(t, expRef, actRef)
		}
	}
}

func AssertUnitRef(t *testing.T, expected, actual *ast.UsesClauseItem) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	assert.Equal(t, expected.Path, actual.Path)
}
