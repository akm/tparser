package pcontext

import (
	"path/filepath"

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
	astcore.DeclarationMap
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

func (c *ProgramContext) IsUnitIdentifier(t *token.Token) bool {
	s := t.Value()
	return c.unitIdentifiers.Include(s) || IsUnitDeclaration(c.DeclarationMap.Get(s)) || c.Units.ByName(s) != nil
}

func (c *ProgramContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
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
