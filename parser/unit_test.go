package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestQualIdInCompoundStmt(t *testing.T) {
	procBar := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Bar"),
		},
	}

	unitFooDeclMap := astcore.NewDeclarationMap()
	unitFooDeclMap.SetDecl(procBar)

	unitFoo := &ast.Unit{
		Ident: &ast.Ident{Name: "foo"},
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				procBar,
			},
		},
		DeclarationMap: unitFooDeclMap,
	}

	declMap := astcore.NewDeclarationMap()
	declMap.SetDecl(unitFoo)

	run := func(name string, text []rune, expected *ast.CompoundStmt) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text, NewContext(declMap, ast.Units{unitFoo}))
			parser.NextToken()
			res, err := parser.ParseCompoundStmt(true)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run("simple call statement with QualId",
		[]rune(`
begin
	foo.Bar;
end;`),
		&ast.CompoundStmt{
			StmtList: ast.StmtList{
				{
					Body: &ast.CallStatement{
						Designator: asttest.NewDesignator(
							&ast.QualId{
								UnitId: &ast.IdentRef{
									Ident: asttest.NewIdent(unitFoo.Ident.Name, asttest.NewIdentLocation(2, 3, 8, 6)),
									Ref:   unitFoo.ToDeclarations()[0],
								},
								Ident: &ast.IdentRef{
									Ident: asttest.NewIdent("Bar", asttest.NewIdentLocation(2, 7, 12, 10)),
									Ref:   procBar.ToDeclarations()[0],
								},
							},
						),
					},
				},
			},
		},
	)
}
