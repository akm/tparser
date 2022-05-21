package pcontext

import (
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type UnitContext struct {
	Parent          *ProgramContext
	Path            string
	unitIdentifiers ext.Strings // TO BE REMOVED
	astcore.DeclMap
}

func NewUnitContext(parent *ProgramContext, args ...interface{}) *UnitContext {
	var path string
	var unitIdentifiers ext.Strings
	var declarationMap astcore.DeclMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
		case ext.Strings:
			unitIdentifiers = v
		case astcore.DeclMap:
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
		DeclMap:         declarationMap,
	}
}

func (c *UnitContext) Clone() Context {
	return &UnitContext{
		Parent:          c.Parent,
		Path:            c.Path,
		unitIdentifiers: c.unitIdentifiers,
		DeclMap:         c.DeclMap,
	}
}
func (c *UnitContext) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *UnitContext) IsUnitIdentifier(token *token.Token) bool {
	s := token.Value()
	return c.unitIdentifiers.Include(s) || IsUnitDeclaration(c.DeclMap.Get(s))
}

func (c *UnitContext) GetDeclarationMap() astcore.DeclMap {
	return c.DeclMap
}

func (c *UnitContext) GetPath() string {
	return c.Path
}
