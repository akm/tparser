package pcontext

import (
	"github.com/akm/tparser/ast"
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
func (c *UnitContext) ImportUnitDecls(usesClause ast.UsesClause) error {
	units := ast.Units{}
	parentUnits := c.Parent.Units
	for _, unitItem := range usesClause {
		if u := parentUnits.ByName(unitItem.Ident.Name); u != nil {
			unitItem.Unit = u
			units = append(units, u)
		}
	}
	localMap := astcore.NewDeclarationMap()
	maps := []astcore.DeclMap{localMap, c.DeclMap}
	for _, unit := range units {
		if err := localMap.Set(unit); err != nil {
			return err
		}
		// TODO declMapに追加する順番はこれでOK？
		// 無関係のユニットAとBに、同じ名前の型や変数が定義されていて、USES A, B; となっていた場合
		// コンテキスト上ではどちらが有効になるのかを確認する
		maps = append(maps, unit.DeclarationMap)
	}
	c.DeclMap = astcore.NewCompositeDeclarationMap(maps...)
	return nil
}

func (c *UnitContext) IsUnitIdentifier(token *token.Token) bool {
	s := token.Value()
	return c.unitIdentifiers.Include(s) || IsUsesClauseItem(c.DeclMap.Get(s))
}

func (c *UnitContext) GetDeclarationMap() astcore.DeclMap {
	return c.DeclMap
}

func (c *UnitContext) GetPath() string {
	return c.Path
}
