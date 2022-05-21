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
	return c.unitIdentifiers.Include(s) || IsUnitDeclaration(c.DeclarationMap.Get(s))
}

func (c *UnitContext) GetDeclarationMap() astcore.DeclarationMap {
	return c.DeclarationMap
}

func (c *UnitContext) GetPath() string {
	return c.Path
}
