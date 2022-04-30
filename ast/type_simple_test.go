package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleType(t *testing.T) {
	t.Run("SimpleType", func(t *testing.T) {
		assert.Implements(t, (*SimpleType)(nil), &RealType{})
		assert.Implements(t, (*Type)(nil), &RealType{})
		assert.Implements(t, (*Node)(nil), &RealType{})
	})
	t.Run("OrdinalType extends SimpleType", func(t *testing.T) {
		objects := []OrdinalType{
			&SubrangeType{},
			EnumeratedType{},
			&OrdIdent{},
		}
		for _, obj := range objects {
			assert.Implements(t, (*OrdinalType)(nil), obj)
			assert.Implements(t, (*SimpleType)(nil), obj)
			assert.Implements(t, (*Type)(nil), obj)
			assert.Implements(t, (*Node)(nil), obj)
		}
	})
}
