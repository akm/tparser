package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
)

type StackableContext struct {
	path    *string
	parent  Context
	declMap astcore.DeclMap
}

var _ Context = (*StackableContext)(nil)

func NewStackableContext(parent Context, args ...interface{}) Context {
	return &StackableContext{
		parent:  parent,
		declMap: astcore.NewDeclMap(),
	}
}

func (c *StackableContext) Clone() Context {
	return &ProgramContext{
		DeclMap: c.declMap,
	}
}

func (c *StackableContext) Get(name string) *astcore.Decl {
	if r := c.declMap.Get(name); r != nil {
		return r
	}
	return c.parent.Get(name)
}

func (c *StackableContext) Set(decl astcore.DeclNode) error {
	return c.declMap.Set(decl)
}

func (c *StackableContext) Overwrite(name string, decl *astcore.Decl) {
	c.declMap.Overwrite(name, decl)
}

func (c *StackableContext) GetPath() string {
	if c.path != nil {
		return *c.path
	}
	return c.parent.GetPath()
}

func (c *StackableContext) StackDeclMap() func() {
	var backup astcore.DeclMap
	c.declMap, backup = astcore.NewChainedDeclMap(c.declMap), c.declMap
	return func() {
		c.declMap = backup
	}
}
