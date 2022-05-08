package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewCaseLabel(arg interface{}, extras ...interface{}) *ast.CaseLabel {
	switch len(extras) {
	case 0:
		return ast.NewCaseLabel(NewConstExpr(arg))
	case 1:
		return ast.NewCaseLabel(NewConstExpr(arg), NewConstExpr(extras[0]))
	default:
		panic("too many extras for asttest.NewCaseLabel")
	}
}
