package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestCaseStmt(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Statement) {
		runSatement(t, name, true, text, expected)
	}

	run(
		"without ELSE",
		[]rune(`
case MyColor of
	Red: X := 1;
	Green: X := 2;
	Blue: X := 3;
	Yellow, Orange, Black: X := 0;
end;
`),
		&ast.Statement{
			Body: &ast.CaseStmt{
				Expression: asttest.NewExpression(asttest.NewQualId("MyColor")),
				Selectors: ast.CaseSelectors{
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel("Red")},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("X"),
								Expression: asttest.NewExpression(ast.NewNumber("1")),
							},
						},
					},
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel("Green")},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("X"),
								Expression: asttest.NewExpression(ast.NewNumber("2")),
							},
						},
					},
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel("Blue")},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("X"),
								Expression: asttest.NewExpression(ast.NewNumber("3")),
							},
						},
					},
					{
						Labels: ast.CaseLabels{
							asttest.NewCaseLabel("Yellow"),
							asttest.NewCaseLabel("Orange"),
							asttest.NewCaseLabel("Black"),
						},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("X"),
								Expression: asttest.NewExpression(ast.NewNumber("0")),
							},
						},
					},
				},
			},
		},
	)

	run(
		"with ELSE",
		[]rune(`
case Selection of
	Done: Form1.Close;
	Compute: CalculateTotal(UnitCost, Quantity);
else
	Beep;
end;	
`),
		&ast.Statement{
			Body: &ast.CaseStmt{
				Expression: asttest.NewExpression(asttest.NewQualId("Selection")),
				Selectors: ast.CaseSelectors{
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel("Done")},
						Statement: &ast.Statement{
							Body: &ast.CallStatement{
								Designator: &ast.Designator{
									QualId: asttest.NewQualId("Form1"),
									Items: []ast.DesignatorItem{
										ast.NewDesignatorItemIdent(asttest.NewIdent("Close")),
									},
								},
							},
						},
					},
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel("Compute")},
						Statement: &ast.Statement{
							Body: &ast.CallStatement{
								Designator: &ast.Designator{
									QualId: asttest.NewQualId("CalculateTotal"),
								},
								ExprList: ast.ExprList{
									asttest.NewExpression(asttest.NewQualId("UnitCost")),
									asttest.NewExpression(asttest.NewQualId("Quantity")),
								},
							},
						},
					},
				},
				Else: ast.StmtList{
					{
						Body: &ast.CallStatement{
							Designator: asttest.NewDesignator("Beep"),
						},
					},
				},
			},
		},
	)

	run(
		"with Subranges",
		[]rune(`
case I of
	1..5: Caption := 'Low';
	6..9: Caption := 'High';
	0, 10..99: Caption := 'Out of range';
else
	Caption := ''; 
end;
`),
		&ast.Statement{
			Body: &ast.CaseStmt{
				Expression: asttest.NewExpression(asttest.NewQualId("I")),
				Selectors: ast.CaseSelectors{
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel(
							asttest.NewNumber("1"),
							asttest.NewNumber("5"),
						)},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("Caption"),
								Expression: asttest.NewExpression(ast.NewString("'Low'")),
							},
						},
					},
					{
						Labels: ast.CaseLabels{asttest.NewCaseLabel(
							asttest.NewNumber("6"),
							asttest.NewNumber("9"),
						)},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("Caption"),
								Expression: asttest.NewExpression(ast.NewString("'High'")),
							},
						},
					},
					{
						Labels: ast.CaseLabels{
							asttest.NewCaseLabel(asttest.NewNumber("0")),
							asttest.NewCaseLabel(
								asttest.NewNumber("10"),
								asttest.NewNumber("99"),
							),
						},
						Statement: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("Caption"),
								Expression: asttest.NewExpression(ast.NewString("'Out of range'")),
							},
						},
					},
				},
				Else: ast.StmtList{
					{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator("Caption"),
							Expression: asttest.NewExpression(ast.NewString("''")),
						},
					},
				},
			},
		},
	)

}
