package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
)

func IsUnitDeclaration(decl *astcore.Declaration) bool {
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.Unit)
	return ok
}
