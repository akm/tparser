package asttest

import "github.com/akm/tparser/ast"

func NewVarDeclAbsoluteIdent(arg interface{}) *ast.VarDeclAbsoluteIdent {
	switch v := arg.(type) {
	case string:
		return ast.NewVarDeclAbsoluteIdent(NewIdent(v))
	default:
		return ast.NewVarDeclAbsoluteIdent(arg)
	}
}
