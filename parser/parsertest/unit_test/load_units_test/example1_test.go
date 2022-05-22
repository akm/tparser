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
			Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(5, 11, 36, 14)),
		},
	}
	declGet := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:       ast.FtFunction,
			Ident:      asttest.NewIdent("Get", asttest.NewIdentLocation(6, 10, 51, 13)),
			ReturnType: asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(6, 15, 56, 22))),
		},
	}

	declCount := &ast.VarDecl{
		IdentList: asttest.NewIdentList(
			asttest.NewIdent("Count", asttest.NewIdentLocation(10, 5, 90, 10)),
		),
		Type:      asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(10, 12, 97, 19))),
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
						Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(12, 11, 123, 14)),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Count", asttest.NewIdentLocation(14, 4, 139, 9)),
												astcore.NewDeclaration(declCount.IdentList[0], declCount),
											),
										),
										Expression: asttest.NewExpression(
											&ast.SimpleExpression{
												Term: &ast.Term{Factor: asttest.NewDesignatorFactor(
													asttest.NewQualId(
														asttest.NewIdent("Count", asttest.NewIdentLocation(14, 13, 148, 18)),
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
						Ident:      asttest.NewIdent("Get", asttest.NewIdentLocation(17, 10, 214, 13)),
						ReturnType: asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(17, 15, 219, 22))),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Result", asttest.NewIdentLocation(19, 4, 239, 10)),
											),
										),
										Expression: asttest.NewExpression(
											asttest.NewQualId(
												asttest.NewIdent("Count",
													asttest.NewIdentLocation(19, 14, 249, 19)),
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

	actualUnitBar.DeclarationMap = nil
	t.Run("bar.pas", func(t *testing.T) {
		asttest.ClearUsesItemsUnit(t, actualUnitBar)
		if !assert.Equal(t, expectedUnitBar, actualUnitBar) {
			asttest.AssertUnit(t, expectedUnitBar, actualUnitBar)
		}
	})

	declProcess := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(5, 11, 36, 18)),
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
				{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(9, 6, 71, 14))},
				{Ident: asttest.NewIdent("bar", asttest.NewIdentLocation(9, 16, 81, 19))},
			},
			DeclSections: ast.DeclSections{
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:  ast.FtProcedure,
						Ident: asttest.NewIdent("Process", asttest.NewIdentLocation(11, 11, 99, 18)),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdentRef("Inc", asttest.NewIdentLocation(13, 4, 119, 7), astcore.NewDeclaration(declInc.Ident, declInc)),
											),
										),
									},
								},
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Writeln", asttest.NewIdentLocation(14, 4, 128, 11)),
											),
										),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												&ast.DesignatorFactor{
													Designator: asttest.NewDesignator(
														&ast.QualId{
															UnitId: asttest.NewIdentRef("bar", asttest.NewIdentLocation(14, 13, 137, 16), astcore.NewDeclaration(expectedUnitBar.Ident, expectedUnitBar)),
															Ident:  asttest.NewIdentRef("Get", asttest.NewIdentLocation(14, 17, 141, 20), astcore.NewDeclaration(declGet.Ident, declGet)),
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
												UnitId: asttest.NewIdentRef("bar", asttest.NewIdentLocation(15, 4, 152, 7), astcore.NewDeclaration(expectedUnitBar.Ident, expectedUnitBar)),
												Ident:  asttest.NewIdentRef("Inc", asttest.NewIdentLocation(15, 8, 156, 11), astcore.NewDeclaration(declInc.Ident, declInc)),
											},
										),
									},
								},
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Writeln", asttest.NewIdentLocation(16, 4, 165, 11)),
											),
										),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												&ast.DesignatorFactor{
													Designator: asttest.NewDesignator(
														asttest.NewQualId(
															asttest.NewIdentRef("bar", asttest.NewIdentLocation(16, 13, 174, 16), astcore.NewDeclaration(expectedUnitBar.Ident, expectedUnitBar)),
															asttest.NewIdentRef("Get", asttest.NewIdentLocation(16, 17, 178, 20), astcore.NewDeclaration(declGet.Ident, declGet)),
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

	actualUnitFoo.DeclarationMap = nil
	t.Run("foo.pas", func(t *testing.T) {
		asttest.ClearUsesItemsUnit(t, actualUnitFoo)
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
					{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(4, 3, 47, 11))},
					{
						Ident: asttest.NewIdent("foo", asttest.NewIdentLocation(5, 3, 60, 6)),
						Path:  ext.StringPtr("'foo.pas'"),
					},
					{
						Ident: asttest.NewIdent("bar", asttest.NewIdentLocation(6, 3, 81, 6)),
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
											UnitId: asttest.NewIdentRef("foo", asttest.NewIdentLocation(9, 4, 120, 7), astcore.NewDeclaration(expectedUnitFoo.Ident, expectedUnitFoo)),
											Ident:  asttest.NewIdentRef("Process", asttest.NewIdentLocation(9, 8, 124, 15), astcore.NewDeclaration(declProcess.Ident, declProcess)),
										},
									),
								},
							},
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										asttest.NewQualId(
											asttest.NewIdent("Readln", asttest.NewIdentLocation(10, 4, 137, 10)),
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
		asttest.ClearUsesItemsUnit(t, actualProg)
		if !assert.Equal(t, expectedProg, actualProg) {
			asttest.AssertProgram(t, expectedProg.Program, actualProg.Program)
		}
	})

}
