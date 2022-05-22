package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
		switch v := args[0].(type) {
		case *ast.IdentRef:
			return ast.NewQualId(nil, v)
		case *ast.Ident:
			return NewQualId(ast.NewIdentRef(v, nil))
		default:
			return NewQualId(NewIdent(args[0]))
		}
	case 2:
		if declaration, ok := args[1].(*astcore.Decl); ok {
			ident := NewIdent(args[0])
			return ast.NewQualId(nil, ast.NewIdentRef(ident, declaration))
		}
		var unitId *ast.IdentRef
		switch v := args[0].(type) {
		case *ast.IdentRef:
			unitId = v
		default:
			unitId = ast.NewIdentRef(NewIdent(args[0]), nil)
		}
		var ident *ast.IdentRef
		switch v := args[1].(type) {
		case *ast.IdentRef:
			ident = v
		default:
			ident = ast.NewIdentRef(NewIdent(args[1]), nil)
		}
		return ast.NewQualId(unitId, ident)
	// case 3:
	// 	declaration, ok := args[2].(*astcore.Declaration)
	// 	if !ok {
	// 		panic(errors.Errorf("unexpected type of args[2] for NewQualId %+v", args))
	// 	}
	// 	r := NewQualId(args[0], args[1])
	// 	r.Ref = declaration
	// 	return r
	default:
		panic(errors.Errorf("unexpected args for NewQualId: %+v", args))
	}
}

func ClearUnitDeclarationMap(u *ast.Unit) {
	u.DeclarationMap = nil
}

func ClearUnitDeclarationMaps(t *testing.T, node ast.Node) {
	err := astcore.WalkDown(node, func(n ast.Node) error {
		switch v := n.(type) {
		case *ast.Unit:
			ClearUnitDeclarationMap(v)
		}
		return nil
	})
	assert.NoError(t, err)
}
