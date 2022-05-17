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
	astcore.DeclarationMap
}

type ContextImpl struct {
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
	return &ContextImpl{
		unitIdentifiers: unitIdentifiers,
		Units:           units,
		DeclarationMap:  declarationMap,
	}
}

func (c *ContextImpl) Clone() Context {
	return &ContextImpl{
		unitIdentifiers: c.unitIdentifiers,
		Units:           c.Units,
		DeclarationMap:  c.DeclarationMap,
	}
}

func (c *ContextImpl) AddUnitIdentifiers(names ...string) {
	c.unitIdentifiers = append(c.unitIdentifiers, names...)
}

func (c *ContextImpl) IsUnitIdentifier(token *token.Token) bool {
	return c.unitIdentifiers.Include(token.Value()) || c.Units.ByName(token.Value()) != nil
}
