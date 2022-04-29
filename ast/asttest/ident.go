package asttest

import (
	"github.com/akm/tparser/ast"
	"github.com/pkg/errors"
)

func NewIdent(arg interface{}) *ast.Ident {
	switch v := arg.(type) {
	case string:
		return &ast.Ident{Name: v}
	case *string:
		return &ast.Ident{Name: *v}
	default:
		return ast.NewIdentFrom(arg)
	}
}

func NewIdentList(args ...interface{}) ast.IdentList {
	switch len(args) {
	case 0:
		panic(errors.Errorf("no arguments are given for NewIdentList"))
	case 1:
		arg := args[0]
		switch v := arg.(type) {
		case string:
			return ast.NewIdentList(NewIdent(v))
		case *string:
			return ast.NewIdentList(NewIdent(v))
		case []string:
			r := make(ast.IdentList, len(v))
			for idx, i := range v {
				r[idx] = NewIdent(i)
			}
			return r
		default:
			return ast.NewIdentList(arg)
		}
	default:
		r := make(ast.IdentList, len(args))
		for i, arg := range args {
			r[i] = NewIdent(arg)
		}
		return r
	}
}
