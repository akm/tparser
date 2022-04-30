package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionHeading(t *testing.T) {
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &ExportedHeading{})
		assert.Implements(t, (*Node)(nil), &FunctionHeading{})
		assert.Implements(t, (*Node)(nil), FormalParameters{})
		assert.Implements(t, (*Node)(nil), &FormalParm{})
		assert.Implements(t, (*Node)(nil), &ParameterType{})
		assert.Implements(t, (*Node)(nil), &Parameter{})
	})
}
