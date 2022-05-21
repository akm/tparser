package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type ProgramContext struct {
	Path            string
	unitIdentifiers ext.Strings // TO BE REMOVED
	Units           ast.Units
	astcore.DeclMap
}

func NewProgramContext(args ...interface{}) *ProgramContext {
	var path string
	var unitIdentifiers ext.Strings
	var units ast.Units
	var declarationMap astcore.DeclMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
		case ext.Strings:
			unitIdentifiers = v
		case ast.Units:
			units = v
		case astcore.DeclMap:
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
		DeclMap:         declarationMap,
	}
}

func (c *ProgramContext) Clone() Context {
	return &ProgramContext{
		Path:            c.Path,
		unitIdentifiers: c.unitIdentifiers,
		Units:           c.Units,
		DeclMap:         c.DeclMap,
	}
}

func (c *ProgramContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *ProgramContext) IsUnitIdentifier(t *token.Token) bool {
	s := t.Value()
	return c.unitIdentifiers.Include(s) || IsUnitDeclaration(c.DeclMap.Get(s)) || c.Units.ByName(s) != nil
}

func (c *ProgramContext) GetDeclarationMap() astcore.DeclMap {
	return c.DeclMap
}

func (c *ProgramContext) GetPath() string {
	return c.Path
}

func (c *ProgramContext) AddUnit(unit *ast.Unit) {
	c.Units = append(c.Units, unit)
}
