package asttest

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

func NewTypeId(unitIdOrIdent interface{}, args ...interface{}) *ast.TypeId {
	var ref *astcore.Declaration
	if len(args) > 0 {
		if v, ok := args[len(args)-1].(*astcore.Declaration); ok {
			ref = v
			args = args[:len(args)-1]
		}
	}
	if len(args) == 0 {
		callArgs := []interface{}{}
		if ref != nil {
			callArgs = append(callArgs, ref)
		}
		return ast.NewTypeId(NewIdent(unitIdOrIdent), callArgs...)
	} else if len(args) == 1 {
		callArgs := []interface{}{*NewIdent(args[0])}
		if ref != nil {
			callArgs = append(callArgs, ref)
		}
		return ast.NewTypeId(NewUnitId(unitIdOrIdent), callArgs...)
	} else {
		panic(errors.Errorf("too many arguments for NewTypeId: %v, %v", unitIdOrIdent, args))
	}
}
