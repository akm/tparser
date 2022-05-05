package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

type Context struct {
	unitIdentifiers ext.Strings // TO BE REMOVED
	Units           ast.Units
	DeclarationMap  astcore.DeclarationMap
}

func NewContext(args ...interface{}) *Context {
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
	return &Context{
		unitIdentifiers: unitIdentifiers,
		Units:           units,
		DeclarationMap:  declarationMap,
	}
}

func (c *Context) IsUnitIdentifier(token *token.Token) bool {
	return c.unitIdentifiers.Include(token.Value()) || c.Units.ByName(token.Value()) != nil
}
