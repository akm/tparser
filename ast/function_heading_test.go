package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionHeading(t *testing.T) {
	t.Run("FunctionHeading", func(t *testing.T) {
		assert.Implements(t, (*CodeBlock)(nil), &FunctionHeading{})
	})
}
