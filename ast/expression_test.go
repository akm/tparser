package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	t.Run("DesignatorItem", func(t *testing.T) {
		assert.Implements(t, (*DesignatorItem)(nil), &DesignatorItemIdent{})
		assert.Implements(t, (*DesignatorItem)(nil), DesignatorItemExprList{})
	})
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &Expression{})
		assert.Implements(t, (*Node)(nil), ExprList{})
		assert.Implements(t, (*Node)(nil), &RelOpSimpleExpression{})
		assert.Implements(t, (*Node)(nil), &SimpleExpression{})
		assert.Implements(t, (*Node)(nil), &AddOpTerm{})
		assert.Implements(t, (*Node)(nil), &Term{})
		assert.Implements(t, (*Node)(nil), &MulOpFactor{})
		assert.Implements(t, (*Node)(nil), &Designator{})
		assert.Implements(t, (*Node)(nil), &DesignatorItemIdent{})
		assert.Implements(t, (*Node)(nil), &DesignatorItemExprList{})
		assert.Implements(t, (*Node)(nil), &SetElement{})
	})
}
