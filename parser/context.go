package parser

import (
	"path/filepath"

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
	ResolvePath(path string) string
	AddUnit(unit *ast.Unit)
	astcore.DeclarationMap
}

// ProgramContext
type ProgramContext struct {
	Path            string
	unitIdentifiers ext.Strings // TO BE REMOVED
	Units           ast.Units
	astcore.DeclarationMap
}

func NewContext(args ...interface{}) Context {
	return NewProgramContext(args...)
}

func NewProgramContext(args ...interface{}) *ProgramContext {
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
			panic(errors.Errorf("unexpected type %T (%v) is given for NewProjectContext", arg, arg))
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
	return &ProgramContext{
		Path:            path,
		unitIdentifiers: unitIdentifiers,
		Units:           units,
		DeclarationMap:  declarationMap,
	}
}

func (c *ProgramContext) Clone() Context {
	return &ProgramContext{
		Path:            c.Path,
		unitIdentifiers: c.unitIdentifiers,
		Units:           c.Units,
		DeclarationMap:  c.DeclarationMap,
	}
}

func (c *ProgramContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *ProgramContext) IsUnitIdentifier(token *token.Token) bool {
	s := token.Value()
	return c.unitIdentifiers.Include(s) || isUnitDeclaration(c.DeclarationMap.Get(s)) || c.Units.ByName(s) != nil
}

func (c *ProgramContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
}

func isUnitDeclaration(decl *astcore.Declaration) bool {
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.Unit)
	return ok
}

func (c *ProgramContext) GetPath() string {
	return c.Path
}

func (c *ProgramContext) SetPath(path string) {
	c.Path = path
}

func (c *ProgramContext) ResolvePath(path string) string {
	dir := filepath.Dir(c.GetPath())
	return filepath.Join(dir, path)
}

func (c *ProgramContext) AddUnit(unit *ast.Unit) {
	c.Units = append(c.Units, unit)
}
func (c *ProgramContext) GetUnits() ast.Units {
	return c.Units
}

// UnitContext

type UnitContext struct {
	Parent          *ProgramContext
	Path            string
	unitIdentifiers ext.Strings // TO BE REMOVED
	astcore.DeclarationMap
}

func NewUnitContext(parent *ProgramContext, args ...interface{}) *UnitContext {
	var path string
	var unitIdentifiers ext.Strings
	var declarationMap astcore.DeclarationMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
		case ext.Strings:
			unitIdentifiers = v
		case astcore.DeclarationMap:
			declarationMap = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewUnitContext", arg, arg))
		}
	}
	if unitIdentifiers == nil {
		unitIdentifiers = ext.Strings{}
	}
	if declarationMap == nil {
		declarationMap = astcore.NewDeclarationMap()
	}
	return &UnitContext{
		Parent:          parent,
		Path:            path,
		unitIdentifiers: unitIdentifiers,
		DeclarationMap:  declarationMap,
	}
}

func (c *UnitContext) Clone() Context {
	return &UnitContext{
		Parent:          c.Parent,
		Path:            c.Path,
		unitIdentifiers: c.unitIdentifiers,
		DeclarationMap:  c.DeclarationMap,
	}
}
func (c *UnitContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *UnitContext) IsUnitIdentifier(token *token.Token) bool {
	s := token.Value()
	return c.unitIdentifiers.Include(s) || isUnitDeclaration(c.DeclarationMap.Get(s))
}

func (c *UnitContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
}

func (c *UnitContext) GetPath() string {
	return c.Path
}

func (c *UnitContext) SetPath(path string) {
	c.Path = path
}

func (c *UnitContext) ResolvePath(path string) string {
	dir := filepath.Dir(c.GetPath())
	return filepath.Join(dir, path)
}

func (c *UnitContext) AddUnit(unit *ast.Unit) {
	panic(errors.Errorf("not implemented"))
}
func (c *UnitContext) GetUnits() ast.Units {
	panic(errors.Errorf("not implemented"))
}

// StackableContext
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

func (c *StackableContext) ResolvePath(path string) string {
	dir := filepath.Dir(c.GetPath())
	return filepath.Join(dir, path)
}

func (c *StackableContext) SetPath(path string) {
	c.path = &path
}

func (c *StackableContext) AddUnit(unit *ast.Unit) {
	panic(errors.Errorf("unexpected call of AddUnit"))
}
