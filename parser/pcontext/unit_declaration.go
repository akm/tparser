package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
)

func IsUsesClauseItem(decl *astcore.Decl) bool {
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.UsesClauseItem)
	return ok
}
