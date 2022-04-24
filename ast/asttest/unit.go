package asttest

import "github.com/akm/tparser/ast"

func NewUnitId(name interface{}) *ast.UnitId {
	switch v := name.(type) {
	case string:
		return ast.NewUnitId(NewIdent(v))
	default:
		return ast.NewUnitId(name)
	}
}
