package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/pkg/errors"
)

type UnitContext struct {
	Parent *ProgramContext
	Path   string
	astcore.DeclMap
}

func NewUnitContext(parent *ProgramContext, args ...interface{}) *UnitContext {
	var path string
	var declarationMap astcore.DeclMap
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			path = v
		case astcore.DeclMap:
			declarationMap = v
		default:
			panic(errors.Errorf("unexpected type %T (%v) is given for NewUnitContext", arg, arg))
		}
	}
	if declarationMap == nil {
		declarationMap = astcore.NewDeclMap()
	}
	return &UnitContext{
		Parent:  parent,
		Path:    path,
		DeclMap: declarationMap,
	}
}

func (c *UnitContext) Clone() Context {
	return &UnitContext{
		Parent:  c.Parent,
		Path:    c.Path,
		DeclMap: c.DeclMap,
	}
}

func (c *UnitContext) ImportUnitDecls(usesClause ast.UsesClause) {
	c.AssignUnits(usesClause)
	units := usesClause.Units().Compact()
	// NewCompositeDeclMap は 先頭から最後に向かって検索するので、mapsにはその順序でDeclMapを追加する
	localMap := astcore.NewDeclMap()
	maps := []astcore.DeclMap{localMap}
	maps = append(maps, units.DeclMaps().Reverse()...)
	maps = append(maps, c.DeclMap)
	c.DeclMap = astcore.NewCompositeDeclMap(maps...)
}

func (c *UnitContext) AssignUnits(usesClause ast.UsesClause) {
	parentUnits := c.Parent.Units
	for _, unitItem := range usesClause {
		if u := parentUnits.ByName(unitItem.Ident.Name); u != nil {
			unitItem.Unit = u
		}
	}
}

func (c *UnitContext) GetPath() string {
	return c.Path
}

func (c *UnitContext) StackDeclMap() func() {
	var backup astcore.DeclMap
	c.DeclMap, backup = astcore.NewChainedDeclMap(c.DeclMap), c.DeclMap
	return func() {
		c.DeclMap = backup
	}
}
