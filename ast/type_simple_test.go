package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleType(t *testing.T) {
	t.Run("OrdinalType extends SimpleType", func(t *testing.T) {
		obj := EnumeratedType{}
		assert.Implements(t, (*OrdinalType)(nil), obj)
		assert.Implements(t, (*SimpleType)(nil), obj)
		assert.Implements(t, (*Type)(nil), obj)
		assert.Implements(t, (*Node)(nil), obj)
	})
}
