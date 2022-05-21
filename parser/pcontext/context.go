package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
)

type Context interface {
	Clone() Context
	AddUnitIdentifiers(names ...string)
	IsUnitIdentifier(token *token.Token) bool
	GetDeclarationMap() astcore.DeclarationMap
	GetPath() string
	astcore.DeclarationMap
}
