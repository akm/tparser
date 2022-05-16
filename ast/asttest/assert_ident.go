package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/stretchr/testify/assert"
)

func AssertIdentList(t *testing.T, expected, actual ast.IdentList) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertIdent(t, exp, act)
		}
	}
}

func AssertIdent(t *testing.T, expected, actual *ast.Ident) {
	if expected != nil {
		if assert.NotNil(t, actual) {
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.Location, actual.Location)
		}
	} else {
		assert.Nil(t, actual)
	}
}

func AssertLocation(t *testing.T, expected, actual *astcore.Location) {
	assert.Equal(t, expected.Path, actual.Path)
	if !assert.Equal(t, expected.Start, actual.Start) {
		AssertPosition(t, expected.Start, actual.Start)
	}
	if !assert.Equal(t, expected.End, actual.End) {
		AssertPosition(t, expected.End, actual.End)
	}
}

func AssertPosition(t *testing.T, expected, actual *astcore.Position) {
	// DO nothing
}

func AssertDeclaration(t *testing.T, expected, actual *astcore.Declaration) {
	if expected != nil {
		if assert.NotNil(t, actual) {
			if !assert.Equal(t, expected.Ident, actual.Ident) {
				AssertIdent(t, expected.Ident, actual.Ident)
			}
			if !assert.Equal(t, expected.Node, actual.Node) {
				// TODO? implement AssertNode if necessary
				// AssertNode(t, expected.Node, actual.Node)
			}
		}
	} else {
		assert.Nil(t, actual)
	}
}
