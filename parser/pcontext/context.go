package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
)

type Context interface {
	Clone() Context
	IsUnitIdentifier(token *token.Token) bool
	GetDeclarationMap() astcore.DeclMap
	GetPath() string
	astcore.DeclMap
}
