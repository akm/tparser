package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewStringType(name interface{}) *ast.TypeId {
	return ast.NewStringType(NewIdent(name))
}

func NewFixedStringType(name interface{}, length *ast.ConstExpr) *ast.FixedStringType {
	return ast.NewFixedStringType(NewIdent(name), length)
}
