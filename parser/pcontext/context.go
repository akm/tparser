package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
)

type Context interface {
	Clone() Context
	GetPath() string
	StackDeclMap() func()
	astcore.DeclMap
}
