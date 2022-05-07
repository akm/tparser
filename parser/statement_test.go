package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestStatement(t *testing.T) {
	run := func(name string, clearLocations bool, text []rune, expected *ast.Program) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProgram()
			if assert.NoError(t, err) {
				if clearLocations {
					asttest.ClearLocations(t, res)
				}
				assert.Equal(t, expected, res)
				assert.Equal(t, expected.ProgramBlock.Block, res.ProgramBlock.Block)
			}

		})
	}

	// TODO implement test after class support
	// 	run(
	// 		"inherited statement", false,
	// 		[]rune(`PROGRAM inherited;
	// begin
	// end.`),
	// 		&ast.Program{
	// 			Ident: asttest.NewIdent("inherited", asttest.NewIdentLocation(1, 9, 8, 14)),
	// 			ProgramBlock: &ast.ProgramBlock{
	// 				Block: &ast.Block{
	// 					CompoundStmt: &ast.CompoundStmt{
	// 						StmtList: ast.StmtList{},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	)

	run(
		"Goto Loop", true,
		[]rune(`PROGRAM GotoLoop;
label START1;
begin
	START1: writeln('hello, world');
	GOTO START1;
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("GotoLoop"),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {
					start1 := asttest.NewLabelId("START1")
					start1Decl := &ast.LabelDeclSection{LabelId: start1}

					return &ast.Block{
						DeclSections: ast.DeclSections{start1Decl},
						CompoundStmt: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									LabelId: start1,
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(asttest.NewIdent("writeln")),
										ExprList: ast.ExprList{
											asttest.NewExpression(asttest.NewString("'hello, world'")),
										},
									},
								},
								&ast.Statement{
									Body: &ast.GotoStatement{
										LabelId: asttest.NewLabelId("START1"),
										Ref:     start1Decl.ToDeclarations()[0],
									},
								},
							},
						},
					}
				}(),
			},
		},
	)
}
