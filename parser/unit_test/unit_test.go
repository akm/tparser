package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				asttest.ClearUnitDeclarationMap(res)
				if !assert.Equal(t, expected, res) {
					if !assert.Equal(t, expected.ImplementationSection, res.ImplementationSection) {
						if !assert.Equal(t, expected.ImplementationSection.DeclSections, res.ImplementationSection.DeclSections) {
							if !assert.Equal(t, expected.ImplementationSection.DeclSections[0], res.ImplementationSection.DeclSections[0]) {
								ef := expected.ImplementationSection.DeclSections[0].(*ast.FunctionDecl)
								af := res.ImplementationSection.DeclSections[0].(*ast.FunctionDecl)
								assert.Equal(t, ef.Block, af.Block)
							}
						}
					}
				}
			}
		})
	}

	run(
		"simplest unit",
		[]rune(`UNIT U1;
			interface
			implementation
			end.`),
		&ast.Unit{
			Ident:                 asttest.NewIdent("U1", asttest.NewIdentLocation(1, 6, 5, 8)),
			InterfaceSection:      &ast.InterfaceSection{},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)

	run(
		"CountUp",
		[]rune(`UNIT U1;
interface
procedure CountUp;
implementation

var cnt: integer;
procedure CountUp;
begin
  cnt := cnt + 1;
end;

end.`),
		&ast.Unit{
			Ident: asttest.NewIdent("U1", asttest.NewIdentLocation(1, 6, 5, 8)),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					&ast.ExportedHeading{
						FunctionHeading: &ast.FunctionHeading{
							Type:  ast.FtProcedure,
							Ident: asttest.NewIdent("CountUp", asttest.NewIdentLocation(3, 12, 29, 19)),
						},
					},
				},
			},
			ImplementationSection: func() *ast.ImplementationSection {
				declCnt := &ast.VarDecl{
					IdentList: asttest.NewIdentList(
						asttest.NewIdent("cnt", asttest.NewIdentLocation(6, 6, 58, 9)),
					),
					Type: asttest.NewOrdIdent(asttest.NewIdent("integer", asttest.NewIdentLocation(6, 11, 63, 18))),
				}
				return &ast.ImplementationSection{
					DeclSections: ast.DeclSections{
						ast.VarSection{declCnt},
						&ast.FunctionDecl{
							FunctionHeading: &ast.FunctionHeading{
								Type:  ast.FtProcedure,
								Ident: asttest.NewIdent("CountUp", asttest.NewIdentLocation(7, 12, 82, 19)),
							},
							Block: &ast.Block{
								Body: &ast.CompoundStmt{
									StmtList: ast.StmtList{
										{
											Body: &ast.AssignStatement{
												Designator: asttest.NewDesignator(
													asttest.NewQualId(
														asttest.NewIdent("cnt", asttest.NewIdentLocation(9, 4, 99, 7)),
														astcore.NewDeclaration(declCnt.IdentList[0], declCnt),
													),
												),
												Expression: asttest.NewExpression(
													&ast.SimpleExpression{
														Term: &ast.Term{Factor: asttest.NewDesignatorFactor(
															asttest.NewQualId(
																asttest.NewIdent("cnt", asttest.NewIdentLocation(9, 11, 106, 14)),
																astcore.NewDeclaration(declCnt.IdentList[0], declCnt),
															),
														)},
														AddOpTerms: []*ast.AddOpTerm{
															{AddOp: "+", Term: asttest.NewTerm(asttest.NewNumber("1"))},
														},
													},
												),
											},
										},
									},
								},
							},
						},
					},
				}
			}(),
		},
	)

	run(
		"setup teardown",
		[]rune(`UNIT U1;
interface
procedure Process;
implementation
uses networks;
procedure Process;
begin
  Ping;
end;

initialization
SetupNetwork;

finalization
TeardownNetwork;

end.`),
		&ast.Unit{
			Ident: asttest.NewIdent("U1", asttest.NewIdentLocation(1, 6, 5, 8)),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					&ast.ExportedHeading{
						FunctionHeading: &ast.FunctionHeading{
							Type:  ast.FtProcedure,
							Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(3, 12, 29, 19)),
						},
					},
				},
			},
			ImplementationSection: func() *ast.ImplementationSection {
				return &ast.ImplementationSection{
					UsesClause: ast.UsesClause{
						{Ident: asttest.NewIdent("networks", asttest.NewIdentLocation(5, 7, 58, 15))},
					},
					DeclSections: ast.DeclSections{
						&ast.FunctionDecl{
							FunctionHeading: &ast.FunctionHeading{
								Type:  ast.FtProcedure,
								Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(6, 12, 78, 19)),
							},
							Block: &ast.Block{
								Body: &ast.CompoundStmt{
									StmtList: ast.StmtList{
										{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator(
													asttest.NewQualId(asttest.NewIdent("Ping", asttest.NewIdentLocation(8, 4, 95, 8))),
												),
											},
										},
									},
								},
							},
						},
					},
				}
			}(),
			InitSection: &ast.InitSection{
				InitializationStmts: ast.StmtList{
					{
						Body: &ast.CallStatement{
							Designator: asttest.NewDesignator(
								asttest.NewQualId(asttest.NewIdent("SetupNetwork", asttest.NewIdentLocation(12, 2, 122, 14))),
							),
						},
					},
				},
				FinalizationStmts: ast.StmtList{
					{
						Body: &ast.CallStatement{
							Designator: asttest.NewDesignator(
								asttest.NewQualId(asttest.NewIdent("TeardownNetwork", asttest.NewIdentLocation(15, 2, 150, 17))),
							),
						},
					},
				},
			},
		},
	)

}

func TestInterfaceSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.InterfaceSection) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseInterfaceSectionUses()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}

		})
	}

	// "unit" is loaded already in goal.go
	run(
		"Uses only",
		[]rune(`INTERFACE USES U1,U2,U3;`),
		&ast.InterfaceSection{
			UsesClause: ast.UsesClause{
				asttest.NewUnitRef(asttest.NewIdent("U1", asttest.NewIdentLocation(1, 16, 15, 18))),
				asttest.NewUnitRef(asttest.NewIdent("U2", asttest.NewIdentLocation(1, 19, 18, 21))),
				asttest.NewUnitRef(asttest.NewIdent("U3", asttest.NewIdentLocation(1, 22, 21, 24))),
			},
		},
	)
}
