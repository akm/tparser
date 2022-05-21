package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
)

type StackableContext struct {
	path            *string
	parent          Context
	unitIdentifiers ext.Strings
	declarationMap  astcore.DeclarationMap
}

func NewStackableContext(parent Context, args ...interface{}) Context {
	return &StackableContext{
		parent:          parent,
		unitIdentifiers: ext.Strings{},
		declarationMap:  astcore.NewDeclarationMap(),
	}
}

func (c *StackableContext) Clone() Context {
	return &ProgramContext{
		unitIdentifiers: c.unitIdentifiers,
		DeclarationMap:  c.declarationMap,
	}
}

func (c *StackableContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *StackableContext) IsUnitIdentifier(token *token.Token) bool {
	return c.unitIdentifiers.Include(token.Value()) || c.parent.IsUnitIdentifier(token)
}

func (c *StackableContext) GetDeclarationMap() astcore.DeclarationMap {
	return astcore.NewCompositeDeclarationMap(c.declarationMap, c.parent.GetDeclarationMap())
}

func (c *StackableContext) Get(name string) *astcore.Declaration {
	if r := c.declarationMap.Get(name); r != nil {
		return r
	}
	return c.parent.Get(name)
}

func (c *StackableContext) Set(decl astcore.Decl) {
	c.declarationMap.Set(decl)
}

func (c *StackableContext) GetPath() string {
	if c.path != nil {
		return *c.path
	}
	return c.parent.GetPath()
}
