package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertStringType(t *testing.T, expected, actual *ast.StringType) {
	assert.Equal(t, expected.Name, actual.Name)
	if !assert.Equal(t, expected.Name, actual.Name) {
		AssertConstExpr(t, expected.Length, actual.Length)
	}
}
