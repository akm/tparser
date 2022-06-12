package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewStringType(name interface{}) ast.StringType {
	return ast.NewStringType(NewIdent(name))
}

func NewFixedStringType(name interface{}, length *ast.ConstExpr) *ast.FixedStringType {
	return ast.NewFixedStringType(NewIdent(name), length)
}
