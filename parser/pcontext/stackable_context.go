package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
)

type StackableContext struct {
	path           *string
	parent         Context
	declarationMap astcore.DeclMap
}

func NewStackableContext(parent Context, args ...interface{}) Context {
	return &StackableContext{
		parent:         parent,
		declarationMap: astcore.NewDeclMap(),
	}
}

func (c *StackableContext) Clone() Context {
	return &ProgramContext{
		DeclMap: c.declarationMap,
	}
}

func (c *StackableContext) Get(name string) *astcore.Decl {
	if r := c.declarationMap.Get(name); r != nil {
		return r
	}
	return c.parent.Get(name)
}

func (c *StackableContext) Set(decl astcore.DeclNode) error {
	return c.declarationMap.Set(decl)
}

func (c *StackableContext) GetPath() string {
	if c.path != nil {
		return *c.path
	}
	return c.parent.GetPath()
}

func (c *StackableContext) StackDeclMap() func() {
	var backup astcore.DeclMap
	c.declarationMap, backup = astcore.NewChainedDeclMap(c.declarationMap), c.declarationMap
	return func() {
		c.declarationMap = backup
	}
}
