package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	t.Run("Factor", func(t *testing.T) {
		assert.Implements(t, (*Factor)(nil), &DesignatorFactor{})
		assert.Implements(t, (*Factor)(nil), &Address{})
		assert.Implements(t, (*Factor)(nil), &ValueFactor{})
		assert.Implements(t, (*Factor)(nil), &Number{})
		assert.Implements(t, (*Factor)(nil), &String{})
		assert.Implements(t, (*Factor)(nil), &Nil{})
		assert.Implements(t, (*Factor)(nil), &Parentheses{})
		assert.Implements(t, (*Factor)(nil), &Not{})
		assert.Implements(t, (*Factor)(nil), SetConstructor{})
		assert.Implements(t, (*Factor)(nil), &TypeCast{})
	})
	t.Run("DesignatorItem", func(t *testing.T) {
		assert.Implements(t, (*DesignatorItem)(nil), &DesignatorItemIdent{})
		assert.Implements(t, (*DesignatorItem)(nil), DesignatorItemExprList{})
		assert.Implements(t, (*DesignatorItem)(nil), &DesignatorItemDereference{})
	})
}
