package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewUnitRef(name interface{}, paths ...string) *ast.UnitRef {
	switch v := name.(type) {
	case string:
		return ast.NewUnitRef(NewIdent(v), paths...)
	default:
		return ast.NewUnitRef(name, paths...)
	}
}
