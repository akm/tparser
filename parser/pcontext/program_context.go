package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

type ProgramContext struct {
	Path  string
	Units ast.Units
	astcore.DeclMap
}

var _ Context = (*ProgramContext)(nil)

func NewProgramContext(args ...interface{}) *ProgramContext {
	var path string
	var units ast.Units
	var declarationMap astcore.DeclMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
		case ast.Units:
			units = v
		case astcore.DeclMap:
			declarationMap = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewProjectContext", arg, arg))
		}
	}
	if units == nil {
		units = ast.Units{}
	}
	if declarationMap == nil {
		declarationMap = astcore.NewChainedDeclMap(ast.EmbeddedTypeDeclMap)
	}
	return &ProgramContext{
		Path:    path,
		Units:   units,
		DeclMap: declarationMap,
	}
}

func (c *ProgramContext) Clone() Context {
	return &ProgramContext{
		Path:    c.Path,
		Units:   c.Units,
		DeclMap: c.DeclMap,
	}
}

func (c *ProgramContext) GetPath() string {
	return c.Path
}

func (c *ProgramContext) AddUnit(unit *ast.Unit) {
	c.Units = append(c.Units, unit)
}

func (c *ProgramContext) StackDeclMap() func() {
	var backup astcore.DeclMap
	c.DeclMap, backup = astcore.NewChainedDeclMap(c.DeclMap), c.DeclMap
	return func() {
		c.DeclMap = backup
	}
}
