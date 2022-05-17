package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Context interface {
	Clone() Context
	AddUnitIdentifiers(names ...string)
	IsUnitIdentifier(token *token.Token) bool
	GetDeclarationMap() astcore.DeclarationMap
	astcore.DeclarationMap
}

type ProjectContext struct {
	unitIdentifiers ext.Strings // TO BE REMOVED
	Units           ast.Units
	astcore.DeclarationMap
}

func NewContext(args ...interface{}) Context {
	var unitIdentifiers ext.Strings
	var units ast.Units
	var declarationMap astcore.DeclarationMap
	for _, arg := range args {
		switch v := arg.(type) {
		case ext.Strings:
			unitIdentifiers = v
		case ast.Units:
			units = v
		case astcore.DeclarationMap:
			declarationMap = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewContext", arg, arg))
		}
	}
	if unitIdentifiers == nil {
		unitIdentifiers = ext.Strings{}
	}
	if units == nil {
		units = ast.Units{}
	}
	if declarationMap == nil {
		declarationMap = astcore.NewDeclarationMap()
	}
	return &ProjectContext{
		unitIdentifiers: unitIdentifiers,
		Units:           units,
		DeclarationMap:  declarationMap,
	}
}

func (c *ProjectContext) Clone() Context {
	return &ProjectContext{
		unitIdentifiers: c.unitIdentifiers,
		Units:           c.Units,
		DeclarationMap:  c.DeclarationMap,
	}
}

func (c *ProjectContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *ProjectContext) IsUnitIdentifier(token *token.Token) bool {
	return c.unitIdentifiers.Include(token.Value()) || c.Units.ByName(token.Value()) != nil
}

func (c *ProjectContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
}

type StackableContext struct {
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
	return &ProjectContext{
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
func (c *StackableContext) Set(d *astcore.Declaration) {
	c.declarationMap.Set(d)
}

func (c *StackableContext) SetDecl(decl astcore.Decl) {
	c.declarationMap.SetDecl(decl)
}

func (c *StackableContext) Keys() ext.Strings {
	return append(c.declarationMap.Keys(), c.parent.Keys()...)
}
