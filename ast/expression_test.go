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
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &Expression{})
		assert.Implements(t, (*Node)(nil), ExprList{})
		assert.Implements(t, (*Node)(nil), &RelOpSimpleExpression{})
		assert.Implements(t, (*Node)(nil), &SimpleExpression{})
		assert.Implements(t, (*Node)(nil), &AddOpTerm{})
		assert.Implements(t, (*Node)(nil), &Term{})
		assert.Implements(t, (*Node)(nil), &MulOpFactor{})
		assert.Implements(t, (*Node)(nil), &DesignatorFactor{})
		assert.Implements(t, (*Node)(nil), &Address{})
		assert.Implements(t, (*Node)(nil), &Designator{})
		assert.Implements(t, (*Node)(nil), &DesignatorItemIdent{})
		assert.Implements(t, (*Node)(nil), &DesignatorItemExprList{})
		assert.Implements(t, (*Node)(nil), &DesignatorItemDereference{})
		assert.Implements(t, (*Node)(nil), &ValueFactor{})
		assert.Implements(t, (*Node)(nil), &Number{})
		assert.Implements(t, (*Node)(nil), &String{})
		assert.Implements(t, (*Node)(nil), &Nil{})
		assert.Implements(t, (*Node)(nil), &Parentheses{})
		assert.Implements(t, (*Node)(nil), &Not{})
		assert.Implements(t, (*Node)(nil), &SetConstructor{})
		assert.Implements(t, (*Node)(nil), &SetElement{})
		assert.Implements(t, (*Node)(nil), &TypeCast{})
	})
}
