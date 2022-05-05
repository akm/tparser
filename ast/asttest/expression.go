package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewExpression(arg interface{}) *ast.Expression {
	switch v := arg.(type) {
	case string:
		return ast.NewExpression(NewIdent(v))
	default:
		return ast.NewExpression(arg)
	}
}

func NewSimpleExpression(arg interface{}) *ast.SimpleExpression {
	switch v := arg.(type) {
	case string:
		return ast.NewSimpleExpression(NewIdent(v))
	default:
		return ast.NewSimpleExpression(arg)
	}
}

func NewTerm(arg interface{}) *ast.Term {
	switch v := arg.(type) {
	case string:
		return ast.NewTerm(NewIdent(v))
	default:
		return ast.NewTerm(arg)
	}
}

func NewDesignatorFactor(arg interface{}) *ast.DesignatorFactor {
	switch v := arg.(type) {
	case string:
		return ast.NewDesignatorFactor(NewIdent(v))
	default:
		return ast.NewDesignatorFactor(arg)
	}
}

func NewDesignator(arg interface{}) *ast.Designator {
	switch v := arg.(type) {
	case string:
		return ast.NewDesignator(NewIdent(v))
	default:
		return ast.NewDesignator(arg)
	}
}

func NewDesignatorItemIdent(arg interface{}) *ast.DesignatorItemIdent {
	switch v := arg.(type) {
	case string:
		return ast.NewDesignatorItemIdent(NewIdent(v))
	default:
		return ast.NewDesignatorItemIdent(arg)
	}
}

func NewNumber(v string) *ast.NumberFactor {
	return ast.NewNumber(v)
}

func NewString(v string) *ast.StringFactor {
	return ast.NewString(v)
}

func NewNil() *ast.Nil {
	return ast.NewNil()
}

func NewSetElement(expr *ast.Expression) *ast.SetElement {
	return ast.NewSetElement(expr)
}
