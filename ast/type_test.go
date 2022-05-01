package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	t.Run("Type", func(t *testing.T) {
		assert.Implements(t, (*Type)(nil), &TypeId{})
		// assert.Implements(t, (*Type)(nil), &SimpleType{}) // this check is in type_simple_test.go
		// assert.Implements(t, (*Type)(nil), &StrucType{})
		// assert.Implements(t, (*Type)(nil), &PointerType{})
		assert.Implements(t, (*Type)(nil), &StringType{})
		// assert.Implements(t, (*Type)(nil), &ProcedureType{})
		// assert.Implements(t, (*Type)(nil), &VariantType{})
		// assert.Implements(t, (*Type)(nil), &ClassRefType{})
	})
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), TypeSection{})
		assert.Implements(t, (*Node)(nil), &TypeDecl{})
		assert.Implements(t, (*Node)(nil), &TypeId{})
	})
}
