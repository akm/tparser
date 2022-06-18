package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewVariantType(name interface{}) *ast.TypeId {
	return ast.NewVariantType(NewIdent(name))
}
