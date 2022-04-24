package asttest

import (
	"github.com/akm/tparser/ast"
)

func NewRealType(name interface{}) *ast.RealType {
	switch v := name.(type) {
	case string:
		return ast.NewRealType(*NewIdent(v))
	default:
		return ast.NewRealType(name)
	}
}

func NewOrdIdent(name interface{}) *ast.OrdIdent {
	switch v := name.(type) {
	case string:
		return ast.NewOrdIdent(NewIdent(v))
	default:
		return ast.NewOrdIdent(name)
	}
}
