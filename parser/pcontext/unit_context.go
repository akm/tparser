package pcontext

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/token"
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
func (c *UnitContext) ImportUnitDecls(usesClause ast.UsesClause) error {
	units := ast.Units{}
	parentUnits := c.Parent.Units
	for _, unitItem := range usesClause {
		if u := parentUnits.ByName(unitItem.Ident.Name); u != nil {
			unitItem.Unit = u
			units = append(units, u)
		}
	}
	localMap := astcore.NewDeclMap()
	maps := []astcore.DeclMap{localMap, c.DeclMap}
	for _, unit := range units {

		// TODO declMapに追加する順番はこれでOK？
		// 無関係のユニットAとBに、同じ名前の型や変数が定義されていて、USES A, B; となっていた場合
		// コンテキスト上ではどちらが有効になるのかを確認する
		maps = append(maps, unit.DeclarationMap)
	}
	c.DeclMap = astcore.NewCompositeDeclMap(maps...)
	return nil
}

func (c *UnitContext) IsUnitIdentifier(token *token.Token) bool {
	s := token.Value()
	decl := c.Get(s)
	if decl == nil {
		return false
	}
	_, ok := decl.Node.(*ast.UsesClauseItem)
	// log.Printf("UnitContext.IsUnitIdentifier(%s) decl.Node: %T %+v", s, decl.Node, decl.Node)
	return ok
}

func (c *UnitContext) GetPath() string {
	return c.Path
}
