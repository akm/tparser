package asttest

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

func NewUnitId(name interface{}) *ast.UnitId {
	switch v := name.(type) {
	case string:
		return ast.NewUnitId(NewIdent(v))
	default:
		return ast.NewUnitId(name)
	}
}

func NewQualId(args ...interface{}) *ast.QualId {
	switch len(args) {
	case 0:
		panic(errors.Errorf("unexpected empty args for NewQualId"))
	case 1:
		return ast.NewQualId(nil, NewIdent(args[0]))
	case 2:
		if declaration, ok := args[1].(*astcore.Declaration); ok {
			r := NewQualId(args[0])
			r.Ref = declaration
			return r
		}
		var unitId *ast.UnitId
		switch v := args[0].(type) {
		case *ast.UnitId:
			unitId = v
		case *ast.Ident:
			unitId = ast.NewUnitId(v)
		default:
			unitId = (*ast.UnitId)(NewIdent(args[0]))
		}
		var ident *ast.Ident
		switch v := args[1].(type) {
		case *ast.Ident:
			unitId = ast.NewUnitId(v)
		default:
			unitId = (*ast.UnitId)(NewIdent(args[0]))
		}
		return ast.NewQualId(unitId, ident)
	case 3:
		declaration, ok := args[2].(*astcore.Declaration)
		if !ok {
			panic(errors.Errorf("unexpected type of args[2] for NewQualId %+v", args))
		}
		r := NewQualId(args[0], args[1])
		r.Ref = declaration
		return r
	default:
		panic(errors.Errorf("unexpected args for NewQualId: %+v", args))
	}
}
