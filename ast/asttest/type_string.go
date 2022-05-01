package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewStringType(name string, args ...interface{}) *ast.StringType {
	return ast.NewStringType(name, args...)
}
