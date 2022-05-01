package asttest

import (
	"github.com/akm/tparser/ast"
	"github.com/pkg/errors"
)

func NewTypeId(unitIdOrIdent interface{}, args ...interface{}) *ast.TypeId {
	if len(args) == 0 {
		return ast.NewTypeId(NewIdent(unitIdOrIdent))
	} else if len(args) == 1 {
		return ast.NewTypeId(
			NewUnitId(unitIdOrIdent),
			*NewIdent(args[0]),
		)
	} else {
		panic(errors.Errorf("too many arguments for NewTypeId: %v, %v", unitIdOrIdent, args))
	}
}
