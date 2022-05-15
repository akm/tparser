package loadunits_test

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestLoadExample1Project(t *testing.T) {
	actualProg, err := parser.ParseProgram("example1.dpr")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "example1", actualProg.Ident.Name)

	if !assert.Len(t, actualProg.Units, 2) {
		return
	}

	actualUnitFoo := actualProg.Units[0]
	actualUnitBar := actualProg.Units[1]
	assert.Equal(t, "foo", actualUnitFoo.Ident.Name)
	assert.Equal(t, "bar", actualUnitBar.Ident.Name)

	declInc := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(5, 12, 36, 15)),
		},
	}
	declGet := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:       ast.FtFunction,
			Ident:      asttest.NewIdent("Get", asttest.NewIdentLocation(6, 11, 51, 14)),
			ReturnType: asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(6, 16, 56, 23))),
		},
	}

	declCount := &ast.VarDecl{
		IdentList: asttest.NewIdentList(
			asttest.NewIdent("Count", asttest.NewIdentLocation(10, 6, 90, 11)),
		),
		Type:      asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(10, 13, 97, 20))),
		ConstExpr: asttest.NewConstExpr(asttest.NewNumber("0")),
	}

	expectedUnitBar := &ast.Unit{
		Path:  "subdir1/bar.pas",
		Ident: asttest.NewIdent("bar", asttest.NewIdentLocation(1, 6, 5, 9)),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{declInc, declGet},
		},
		ImplementationSection: &ast.ImplementationSection{
			DeclSections: ast.DeclSections{
				ast.VarSection{declCount},
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:  ast.FtProcedure,
						Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(12, 12, 123, 15)),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Count", asttest.NewIdentLocation(14, 5, 139, 10)),
												astcore.NewDeclaration(declCount.IdentList[0], declCount),
											),
										),
										Expression: asttest.NewExpression(
											&ast.SimpleExpression{
												Term: &ast.Term{Factor: asttest.NewDesignatorFactor(
													asttest.NewQualId(
														asttest.NewIdent("Count", asttest.NewIdentLocation(14, 14, 148, 19)),
														astcore.NewDeclaration(declCount.IdentList[0], declCount),
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
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:       ast.FtFunction,
						Ident:      asttest.NewIdent("Get", asttest.NewIdentLocation(17, 11, 214, 14)),
						ReturnType: asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(17, 16, 219, 23))),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Result", asttest.NewIdentLocation(19, 5, 239, 11)),
											),
										),
										Expression: asttest.NewExpression(
											asttest.NewQualId(
												asttest.NewIdent("Count", asttest.NewIdentLocation(19, 15, 249, 20)),
												astcore.NewDeclaration(declCount.IdentList[0], declCount),
											),
										),
									},
								},
							},
						},
					},
				}},
		},
	}

	t.Run("bar.pas", func(t *testing.T) {
		if !assert.Equal(t, expectedUnitBar, actualUnitBar) {
			asttest.AssertUnit(t, expectedUnitBar, actualUnitBar)
		}
	})

	declProcess := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(5, 12, 36, 19)),
		},
	}

	expectedUnitFoo := &ast.Unit{
		Path:  "foo.pas",
		Ident: asttest.NewIdent("foo", asttest.NewIdentLocation(1, 6, 5, 9)),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{declProcess},
		},
		ImplementationSection: &ast.ImplementationSection{
			UsesClause: ast.UsesClause{
				{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(9, 7, 71, 15))},
				{Ident: asttest.NewIdent("bar", asttest.NewIdentLocation(9, 17, 81, 20))},
			},
			DeclSections: ast.DeclSections{
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:  ast.FtProcedure,
						Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(11, 12, 99, 19)),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Inc", asttest.NewIdentLocation(13, 5, 119, 8)),
												astcore.NewDeclaration(declInc.Ident, declInc),
											),
										),
									},
								},
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Writeln", asttest.NewIdentLocation(14, 5, 128, 12)),
											),
										),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												&ast.DesignatorFactor{
													Designator: asttest.NewDesignator(
														&ast.QualId{
															UnitId: asttest.NewIdentRef("bar", asttest.NewIdentLocation(14, 14, 137, 17)),
															Ident:  asttest.NewIdentRef("Get", asttest.NewIdentLocation(14, 18, 141, 21)),
														},
													),
												},
											),
										},
									},
								},
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											&ast.QualId{
												UnitId: asttest.NewIdentRef("bar", asttest.NewIdentLocation(15, 5, 152, 8)),
												Ident:  asttest.NewIdentRef("Inc", asttest.NewIdentLocation(15, 9, 156, 12), astcore.NewDeclaration(declInc.Ident, declInc)),
											},
										),
									},
								},
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Writeln", asttest.NewIdentLocation(16, 5, 165, 12)),
											),
										),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												&ast.DesignatorFactor{
													Designator: asttest.NewDesignator(
														asttest.NewQualId(
															(*ast.UnitId)(asttest.NewIdent("bar", asttest.NewIdentLocation(16, 14, 174, 17))),
															asttest.NewIdent("Get", asttest.NewIdentLocation(16, 18, 178, 21)),
														),
													),
												},
											),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Run("foo.pas", func(t *testing.T) {
		if !assert.Equal(t, expectedUnitFoo, actualUnitFoo) {
			asttest.AssertUnit(t, expectedUnitFoo, actualUnitFoo)
		}
	})

	expectedProg := &parser.Program{
		Program: &ast.Program{
			Path:  "example1.dpr",
			Ident: asttest.NewIdent("example1", asttest.NewIdentLocation(1, 9, 8, 17)),
			ProgramBlock: &ast.ProgramBlock{
				UsesClause: ast.UsesClause{
					{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(4, 4, 47, 12))},
					{
						Ident: asttest.NewIdent("foo", asttest.NewIdentLocation(5, 4, 60, 7)),
						Path:  ext.StringPtr("'foo.pas'"),
					},
					{
						Ident: asttest.NewIdent("bar", asttest.NewIdentLocation(6, 4, 81, 7)),
						Path:  ext.StringPtr("'subdir1\\bar.pas'"),
					},
				},
				Block: &ast.Block{
					Body: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										&ast.QualId{
											UnitId: asttest.NewIdentRef("foo", asttest.NewIdentLocation(9, 5, 120, 8)),
											Ident:  asttest.NewIdentRef("Process", asttest.NewIdentLocation(9, 9, 124, 16), astcore.NewDeclaration(declProcess.Ident, declProcess)),
										},
									),
								},
							},
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										asttest.NewQualId(
											asttest.NewIdent("Readln", asttest.NewIdentLocation(10, 5, 137, 11)),
										),
									),
								},
							},
						},
					},
				},
			},
		},
		Units: ast.Units{expectedUnitFoo, expectedUnitBar},
	}

	t.Run("example1.dpr", func(t *testing.T) {
		if !assert.Equal(t, expectedProg, actualProg) {
			asttest.AssertProgram(t, expectedProg.Program, actualProg.Program)
		}
	})

}
