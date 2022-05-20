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

func TestLoadSimpleProj1(t *testing.T) {
	t.Skip()
	actualProg, err := parser.ParseProgram("simple_proj1.dpr")
	assert.NoError(t, err)

	actualProg.Program.ProgramBlock.DeclSections = nil
	assert.Equal(t, "simple_proj1", actualProg.Ident.Name)

	if !assert.Len(t, actualProg.Units, 2) {
		return
	}

	actualUnitCallInc := actualProg.Units[0]
	actualUnitCnt := actualProg.Units[1]
	assert.Equal(t, "call_inc", actualUnitCallInc.Ident.Name)
	assert.Equal(t, "cnt", actualUnitCnt.Ident.Name)

	declInc := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(6, 11, 61, 14)),
		},
	}
	declCount := &ast.VarDecl{
		IdentList: asttest.NewIdentList(
			asttest.NewIdent("Count", asttest.NewIdentLocation(5, 5, 30, 10)),
		),
		Type:      asttest.NewOrdIdent(asttest.NewIdent("Integer", asttest.NewIdentLocation(5, 12, 37, 19))),
		ConstExpr: asttest.NewConstExpr(asttest.NewNumber("0")),
	}

	expectedUnitBar := &ast.Unit{
		Path:  "subdir1/cnt.pas",
		Ident: asttest.NewIdent("cnt", asttest.NewIdentLocation(1, 5, 5, 8)),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{ast.VarSection{declCount}, declInc},
		},
		ImplementationSection: &ast.ImplementationSection{
			DeclSections: ast.DeclSections{
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:  ast.FtProcedure,
						Ident: asttest.NewIdent("Inc", asttest.NewIdentLocation(10, 11, 97, 14)),
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
			},
		},
	}

	t.Run("inc.pas", func(t *testing.T) {
		if !assert.Equal(t, expectedUnitBar, actualUnitCnt) {
			asttest.AssertUnit(t, expectedUnitBar, actualUnitCnt)
		}
	})

	declProcess := &ast.ExportedHeading{
		FunctionHeading: &ast.FunctionHeading{
			Type:  ast.FtProcedure,
			Ident: asttest.NewIdent("CallInc", asttest.NewIdentLocation(5, 11, 41, 18)),
		},
	}

	expectedUnitFoo := &ast.Unit{
		Path:  "call_inc.pas",
		Ident: asttest.NewIdent("call_inc", asttest.NewIdentLocation(1, 5, 5, 13)),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{declProcess},
		},
		ImplementationSection: &ast.ImplementationSection{
			UsesClause: ast.UsesClause{
				{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(9, 6, 76, 14))},
				{Ident: asttest.NewIdent("cnt", asttest.NewIdentLocation(9, 16, 86, 19))},
			},
			DeclSections: ast.DeclSections{
				&ast.FunctionDecl{
					FunctionHeading: &ast.FunctionHeading{
						Type:  ast.FtProcedure,
						Ident: asttest.NewIdent("CallInc", asttest.NewIdentLocation(11, 11, 104, 18)),
					},
					Block: &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(
												asttest.NewIdent("Inc", asttest.NewIdentLocation(13, 4, 124, 7)),
												astcore.NewDeclaration(declInc.Ident, declInc),
											),
										),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Run("call_inc.pas", func(t *testing.T) {
		if !assert.Equal(t, expectedUnitFoo, actualUnitCallInc) {
			asttest.AssertUnit(t, expectedUnitFoo, actualUnitCallInc)
		}
	})

	expectedProg := &parser.Program{
		Program: &ast.Program{
			Path:  "simple_proj1.dpr",
			Ident: asttest.NewIdent("simple_proj1", asttest.NewIdentLocation(1, 8, 8, 20)),
			ProgramBlock: &ast.ProgramBlock{
				UsesClause: ast.UsesClause{
					{Ident: asttest.NewIdent("SysUtils", asttest.NewIdentLocation(4, 3, 51, 11))},
					{
						Ident: asttest.NewIdent("call_inc", asttest.NewIdentLocation(5, 3, 64, 11)),
						Path:  ext.StringPtr("'call_inc.pas'"),
					},
					{
						Ident: asttest.NewIdent("cnt", asttest.NewIdentLocation(6, 3, 95, 6)),
						Path:  ext.StringPtr("'subdir1\\cnt.pas'"),
					},
				},
				Block: &ast.Block{
					Body: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										&ast.QualId{
											Ident: asttest.NewIdentRef(
												"CallInc",
												asttest.NewIdentLocation(9, 3, 133, 10),
												astcore.NewDeclaration(declProcess.Ident, declProcess),
											),
										},
									),
								},
							},
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										asttest.NewQualId(
											asttest.NewIdent("Readln", asttest.NewIdentLocation(10, 3, 145, 9)),
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

	t.Run("simple_proj1.dpr", func(t *testing.T) {
		if !assert.Equal(t, expectedProg, actualProg) {
			asttest.AssertProgram(t, expectedProg.Program, actualProg.Program)
		}
	})

}
