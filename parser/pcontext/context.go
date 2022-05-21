package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
)

type Context interface {
	Clone() Context
	AddUnitIdentifiers(names ...string)
	IsUnitIdentifier(token *token.Token) bool
	GetDeclarationMap() astcore.DeclarationMap
	GetPath() string
	SetPath(path string)
	ResolvePath(path string) string
	AddUnit(unit *ast.Unit)
	astcore.DeclarationMap
}
