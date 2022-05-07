package asttest

import "github.com/akm/tparser/ast"

func NewLabelId(arg interface{}) *ast.LabelId {
	switch v := arg.(type) {
	// case *ast.LabelId:
	// 	return v
	case *ast.Ident:
		return ast.NewLabelId(v)
	default:
		return ast.NewLabelId(NewIdent(arg))
	}
}
