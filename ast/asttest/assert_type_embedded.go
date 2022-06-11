package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertTypeEmbedded(t *testing.T, expected, actual *ast.TypeEmbedded) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	assert.Equal(t, expected.Kind, actual.Kind)
}
