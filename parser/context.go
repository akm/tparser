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
	GetPath() string
	SetPath(path string)
	GetUnits() ast.Units
	astcore.DeclarationMap
}

type ProjectContext struct {
	Path            string
	unitIdentifiers ext.Strings // TO BE REMOVED
	Units           ast.Units
	astcore.DeclarationMap
}

func NewContext(args ...interface{}) Context {
	return NewProjectContext(args...)
}

func NewProjectContext(args ...interface{}) *ProjectContext {
	var path string
	var unitIdentifiers ext.Strings
	var units ast.Units
	var declarationMap astcore.DeclarationMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
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
		Path:            path,
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
	s := token.Value()
	return c.unitIdentifiers.Include(s) || isUnitDeclaration(c.DeclarationMap.Get(s)) || c.Units.ByName(s) != nil
}

func (c *ProjectContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
}

func isUnitDeclaration(decl *astcore.Declaration) bool {
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.Unit)
	return ok
}

func (c *ProjectContext) GetPath() string {
	return c.Path
}

func (c *ProjectContext) SetPath(path string) {
	c.Path = path
}

func (c *ProjectContext) GetUnits() ast.Units {
	return c.Units
}

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
func (c *StackableContext) GetPath() string {
	if c.path != nil {
		return *c.path
	}
	return c.parent.GetPath()
}

func (c *StackableContext) SetPath(path string) {
	c.path = &path
}

func (c *StackableContext) GetUnits() ast.Units {
	return c.parent.GetUnits()
}
