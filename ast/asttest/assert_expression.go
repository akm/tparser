package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertExprList(t *testing.T, expected, actual ast.ExprList) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertExpression(t, exp, act)
		}
	}
}

func AssertExpression(t *testing.T, expected, actual *ast.Expression) {
	if !assert.Equal(t, expected.SimpleExpression, actual.SimpleExpression) {
		AssertSimpleExpression(t, expected.SimpleExpression, actual.SimpleExpression)
	}
	if !assert.Equal(t, expected.RelOpSimpleExpressions, actual.RelOpSimpleExpressions) {
		AssertRelOpSimpleExpressions(t, expected.RelOpSimpleExpressions, actual.RelOpSimpleExpressions)
	}
}

func AssertRelOpSimpleExpressions(t *testing.T, expected, actual ast.RelOpSimpleExpressions) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertRelOpSimpleExpression(t, exp, act)
		}
	}
}

func AssertRelOpSimpleExpression(t *testing.T, expected, actual *ast.RelOpSimpleExpression) {
	assert.Equal(t, expected.RelOp, actual.RelOp)
	if !assert.Equal(t, expected.SimpleExpression, actual.SimpleExpression) {
		AssertSimpleExpression(t, expected.SimpleExpression, actual.SimpleExpression)
	}
}

func AssertSimpleExpression(t *testing.T, expected, actual *ast.SimpleExpression) {
	assert.Equal(t, expected.UnaryOp, actual.UnaryOp)
	if !assert.Equal(t, expected.Term, actual.Term) {
		AssertTerm(t, expected.Term, actual.Term)
	}
	if !assert.Equal(t, expected.AddOpTerms, actual.AddOpTerms) {
		AssertAddOpTerms(t, expected.AddOpTerms, actual.AddOpTerms)
	}
}

func AssertAddOpTerms(t *testing.T, expected, actual ast.AddOpTerms) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertAddOpTerm(t, exp, act)
		}
	}
}

func AssertAddOpTerm(t *testing.T, expected, actual *ast.AddOpTerm) {
	assert.Equal(t, expected.AddOp, actual.AddOp)
	if !assert.Equal(t, expected.Term, actual.Term) {
		AssertTerm(t, expected.Term, actual.Term)
	}
}

func AssertTerm(t *testing.T, expected, actual *ast.Term) {
	if !assert.Equal(t, expected.Factor, actual.Factor) {
		AssertFactor(t, expected.Factor, actual.Factor)
	}
	if !assert.Equal(t, expected.MulOpFactors, actual.MulOpFactors) {
		AssertMulOpFactors(t, expected.MulOpFactors, actual.MulOpFactors)
	}
}

func AssertMulOpFactors(t *testing.T, expected, actual ast.MulOpFactors) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertMulOpFactor(t, exp, act)
		}
	}
}

func AssertMulOpFactor(t *testing.T, expected, actual *ast.MulOpFactor) {
	assert.Equal(t, expected.MulOp, actual.MulOp)
	if !assert.Equal(t, expected.Factor, actual.Factor) {
		AssertFactor(t, expected.Factor, actual.Factor)
	}
}

func AssertFactor(t *testing.T, expected, actual ast.Factor) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.DesignatorFactor:
		AssertDesignatorFactor(t, exp, actual.(*ast.DesignatorFactor))
	case *ast.Address:
		AssertAddress(t, exp, actual.(*ast.Address))
	case *ast.NumberFactor:
		AssertNumberFactor(t, exp, actual.(*ast.NumberFactor))
	case *ast.StringFactor:
		AssertStringFactor(t, exp, actual.(*ast.StringFactor))
	case *ast.ValueFactor:
		AssertValueFactor(t, exp, actual.(*ast.ValueFactor))
	case *ast.Nil:
		AssertNil(t, exp, actual.(*ast.Nil))
	case *ast.Parentheses:
		AssertParentheses(t, exp, actual.(*ast.Parentheses))
	case *ast.Not:
		AssertNot(t, exp, actual.(*ast.Not))
	case *ast.SetConstructor:
		AssertSetConstructor(t, exp, actual.(*ast.SetConstructor))
	case *ast.TypeCast:
		AssertTypeCast(t, exp, actual.(*ast.TypeCast))
	default:
		assert.Fail(t, "unexpected Factor type %T", exp)
	}
}

func AssertDesignatorFactor(t *testing.T, expected, actual *ast.DesignatorFactor) {
	if !assert.Equal(t, expected.Designator, actual.Designator) {
		AssertDesignator(t, expected.Designator, actual.Designator)
	}
	if !assert.Equal(t, expected.ExprList, actual.ExprList) {
		AssertExprList(t, expected.ExprList, actual.ExprList)
	}
}

func AssertAddress(t *testing.T, expected, actual *ast.Address) {
	if !assert.Equal(t, expected.Designator, actual.Designator) {
		AssertDesignator(t, expected.Designator, actual.Designator)
	}
}

func AssertDesignator(t *testing.T, expected, actual *ast.Designator) {
	if !assert.Equal(t, expected.QualId, actual.QualId) {
		AssertQualId(t, expected.QualId, actual.QualId)
	}
	if !assert.Equal(t, expected.Items, actual.Items) {
		AssertDesignatorItems(t, expected.Items, actual.Items)
	}
}

func AssertDesignatorItems(t *testing.T, expected, actual ast.DesignatorItems) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertDesignatorItem(t, exp, act)
		}
	}
}

func AssertDesignatorItem(t *testing.T, expected, actual ast.DesignatorItem) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.DesignatorItemIdent:
		AssertDesignatorItemIdent(t, exp, actual.(*ast.DesignatorItemIdent))
	case ast.DesignatorItemExprList:
		AssertDesignatorItemExprList(t, exp, actual.(ast.DesignatorItemExprList))
	case *ast.DesignatorItemDereference:
		AssertDesignatorItemDereference(t, exp, actual.(*ast.DesignatorItemDereference))
	default:
		assert.Fail(t, "unexpected DesignatorItem type %T", exp)
	}
}

func AssertDesignatorItemIdent(t *testing.T, expected, actual *ast.DesignatorItemIdent) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
}

func AssertDesignatorItemExprList(t *testing.T, expected, actual ast.DesignatorItemExprList) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertExpression(t, exp, act)
		}
	}
}

func AssertDesignatorItemDereference(t *testing.T, expected, actual *ast.DesignatorItemDereference) {
	// Do nothing
}

func AssertNumberFactor(t *testing.T, expected, actual *ast.NumberFactor) {
	assert.Equal(t, expected.Value, actual.Value)
}

func AssertStringFactor(t *testing.T, expected, actual *ast.StringFactor) {
	assert.Equal(t, expected.Value, actual.Value)
}

func AssertValueFactor(t *testing.T, expected, actual *ast.ValueFactor) {
	assert.Equal(t, expected.Value, actual.Value)
}

func AssertNil(t *testing.T, expected, actual *ast.Nil) {
	// Do nothing
}

func AssertParentheses(t *testing.T, expected, actual *ast.Parentheses) {
	if !assert.Equal(t, expected.Expression, actual.Expression) {
		AssertExpression(t, expected.Expression, actual.Expression)
	}
}

func AssertNot(t *testing.T, expected, actual *ast.Not) {
	if !assert.Equal(t, expected.Factor, actual.Factor) {
		AssertFactor(t, expected.Factor, actual.Factor)
	}
}

func AssertSetConstructor(t *testing.T, expected, actual *ast.SetConstructor) {
	if !assert.Equal(t, expected.SetElements, actual.SetElements) {
		AssertSetElements(t, expected.SetElements, actual.SetElements)
	}
}

func AssertSetElements(t *testing.T, expected, actual []*ast.SetElement) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertSetElement(t, exp, act)
		}
	}
}

func AssertSetElement(t *testing.T, expected, actual *ast.SetElement) {
	if !assert.Equal(t, expected.Expression, actual.Expression) {
		AssertExpression(t, expected.Expression, actual.Expression)
	}
	if !assert.Equal(t, expected.SubRangeEnd, actual.SubRangeEnd) {
		AssertExpression(t, expected.SubRangeEnd, actual.SubRangeEnd)
	}
}

func AssertTypeCast(t *testing.T, expected, actual *ast.TypeCast) {
	if !assert.Equal(t, expected.TypeId, actual.TypeId) {
		AssertTypeId(t, expected.TypeId, actual.TypeId)
	}
	if !assert.Equal(t, expected.Expression, actual.Expression) {
		AssertExpression(t, expected.Expression, actual.Expression)
	}
}
