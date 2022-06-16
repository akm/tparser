package asttest

import (
	"github.com/akm/tparser/ast"
	"github.com/pkg/errors"
)

func NewRealType(name interface{}) ast.RealType {
	switch v := name.(type) {
	case string:
		return ast.NewRealType(NewIdent(v))
	case *ast.Ident:
		return ast.NewRealType(v)
	default:
		panic(errors.Errorf("invalid type %T for asttest.NewRealType %+v", name, name))
	}
}

func NewOrdIdent(name interface{}) ast.OrdIdent {
	switch v := name.(type) {
	case string:
		return ast.NewOrdIdent(NewIdent(v))
	case *ast.Ident:
		return ast.NewOrdIdent(v)
	default:
		panic(errors.Errorf("invalid type %T for NewOrdIdent %+v", name, name))
	}
}

func NewOrdIdentWithIdent(v *ast.Ident) *ast.TypeId {
	return ast.NewOrdIdent(v)
}
