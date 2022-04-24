package asttest

import "github.com/akm/tparser/ast"

func NewConstExpr(arg interface{}) *ast.ConstExpr {
	return NewExpression(arg)
}
