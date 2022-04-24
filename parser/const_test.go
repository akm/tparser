package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithConstSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				asttest.ClearAllRange(res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"const declaration in unit",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		CONST
		  MaxValue = 237;
		  Message1 = 'Out of memory';
		  Max: Integer = 100;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ConstSection{
						{Ident: ast.Ident("MaxValue"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("237"))},
						{Ident: ast.Ident("Message1"), ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'"))},
						{Ident: ast.Ident("Max"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("100")), Type: &ast.OrdIdent{Name: "Integer"}},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestConstSectionl(t *testing.T) {
	run := func(name string, withRange bool, text []rune, expected ast.ConstSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseConstSection()
			if assert.NoError(t, err) {
				if !withRange {
					asttest.ClearAllRange(res)
				}
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"number const1",
		true,
		[]rune(`CONST MaxValue = 237;`),
		ast.ConstSection{
			{
				Ident:     ast.Ident("MaxValue"),
				ConstExpr: *ast.NewConstExpr(ast.NewNumber("237")),
				CodeBlockNode: ast.CodeBlockNode{
					Range: &ast.CodeRange{
						Path:  "",
						Start: ast.CodePosition{Index: 6, Line: 1, Col: 7},
						End:   ast.CodePosition{Index: 20, Line: 1, Col: 21},
					},
				},
			},
		},
	)
	run(
		"number const2",
		true,
		[]rune(`CONST Max: Integer = 100;`),
		ast.ConstSection{
			{
				Ident:     ast.Ident("Max"),
				ConstExpr: *ast.NewConstExpr(ast.NewNumber("100")), Type: &ast.OrdIdent{Name: "Integer"},
				CodeBlockNode: ast.CodeBlockNode{
					Range: &ast.CodeRange{
						Path:  "",
						Start: ast.CodePosition{Index: 6, Line: 1, Col: 7},
						End:   ast.CodePosition{Index: 24, Line: 1, Col: 25},
					},
				},
			},
		},
	)
	run(
		"message as identifier",
		true,
		[]rune(`CONST Message = 'Out of memory';`),
		ast.ConstSection{
			{
				Ident:     ast.Ident("Message"),
				ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'")),
				CodeBlockNode: ast.CodeBlockNode{
					Range: &ast.CodeRange{
						Path:  "",
						Start: ast.CodePosition{Index: 6, Line: 1, Col: 7},
						End:   ast.CodePosition{Index: 31, Line: 1, Col: 32},
					},
				},
			},
		},
	)

	run(
		"examples in Language Guide",
		false,
		[]rune(`
		const
			Min = 0;
			Max = 100;
			Center = (Max - Min) div 2;
			Beta = Chr(225);
			NumChars = Ord('Z') - Ord('A') + 1;
			Message = 'Out of memory';
			ErrStr = ' Error: ' + Message + '. ';
			ErrPos = 80 - Length(ErrStr) div 2;
			Ln10 = 2.302585092994045684;
			Ln10R = 1 / Ln10;
			Numeric = ['0'..'9'];
			Alpha = ['A'..'Z', 'a'..'z'];
			AlphaNum = Alpha + Numeric;
		`),
		ast.ConstSection{
			{Ident: ast.Ident("Min"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("0"))},
			{Ident: ast.Ident("Max"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("100"))},
			// Center = (Max - Min) div 2;
			{Ident: ast.Ident("Center"), ConstExpr: *ast.NewConstExpr(
				&ast.Term{
					Factor: &ast.Parentheses{
						Expression: ast.Expression{
							SimpleExpression: ast.SimpleExpression{
								Term: *ast.NewTerm("Max"),
								AddOpTerms: []*ast.AddOpTerm{
									{AddOp: "-", Term: *ast.NewTerm("Min")},
								},
							},
						},
					},
					MulOpFactors: []*ast.MulOpFactor{
						{MulOp: "DIV", Factor: ast.NewNumber("2")},
					},
				},
			)},
			// Beta = Chr(225);
			{Ident: ast.Ident("Beta"), ConstExpr: *ast.NewConstExpr(
				&ast.DesignatorFactor{
					Designator: ast.Designator{
						QualId: ast.QualId{Ident: ast.Ident("Chr")},
					},
					ExprList: ast.ExprList{
						ast.NewConstExpr(ast.NewNumber("225")),
					},
				},
			)},
			// NumChars = Ord('Z') - Ord('A') + 1;
			{Ident: ast.Ident("NumChars"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm(
						&ast.DesignatorFactor{
							Designator: ast.Designator{
								QualId: ast.QualId{Ident: ast.Ident("Ord")},
							},
							ExprList: ast.ExprList{
								ast.NewConstExpr(ast.NewString("'Z'")),
							},
						},
					),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "-",
							Term: *ast.NewTerm(
								&ast.DesignatorFactor{
									Designator: ast.Designator{
										QualId: ast.QualId{Ident: ast.Ident("Ord")},
									},
									ExprList: ast.ExprList{
										ast.NewConstExpr(ast.NewString("'A'")),
									},
								},
							),
						},
						{
							AddOp: "+",
							Term:  *ast.NewTerm(ast.NewNumber("1")),
						},
					},
				},
			)},
			// Message = 'Out of memory';
			{Ident: ast.Ident("Message"), ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'"))},
			// ErrStr = ' Error: ' + Message + '. ';
			{Ident: ast.Ident("ErrStr"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm(ast.NewString("' Error: '")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "+",
							Term:  *ast.NewTerm("Message"),
						},
						{
							AddOp: "+",
							Term:  *ast.NewTerm(ast.NewString("'. '")),
						},
					},
				},
			)},
			// ErrPos = 80 - Length(ErrStr) div 2;
			{Ident: ast.Ident("ErrPos"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm(ast.NewNumber("80")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "-",
							Term: ast.Term{
								Factor: &ast.DesignatorFactor{
									Designator: ast.Designator{
										QualId: ast.QualId{Ident: ast.Ident("Length")},
									},
									ExprList: ast.ExprList{
										ast.NewConstExpr("ErrStr"),
									},
								},
								MulOpFactors: []*ast.MulOpFactor{
									{MulOp: "DIV", Factor: ast.NewNumber("2")},
								},
							},
						},
					},
				},
			)},
			// Ln10 = 2.302585092994045684;
			{Ident: ast.Ident("Ln10"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("2.302585092994045684"))},
			// Ln10R = 1 / Ln10;
			{Ident: ast.Ident("Ln10R"), ConstExpr: *ast.NewConstExpr(
				&ast.Term{
					Factor: ast.NewNumber("1"),
					MulOpFactors: []*ast.MulOpFactor{
						{MulOp: "/", Factor: ast.NewDesignatorFactor("Ln10")},
					},
				},
			)},
			// Numeric = ['0'..'9'];
			{Ident: ast.Ident("Numeric"), ConstExpr: *ast.NewConstExpr(
				&ast.SetConstructor{
					SetElements: []*ast.SetElement{
						{
							Expression:  *ast.NewConstExpr(ast.NewString("'0'")),
							SubRangeEnd: ast.NewConstExpr(ast.NewString("'9'")),
						},
					},
				},
			)},
			// Alpha = ['A'..'Z', 'a'..'z'];
			{Ident: ast.Ident("Alpha"), ConstExpr: *ast.NewConstExpr(
				&ast.SetConstructor{
					SetElements: []*ast.SetElement{
						{
							Expression:  *ast.NewConstExpr(ast.NewString("'A'")),
							SubRangeEnd: ast.NewConstExpr(ast.NewString("'Z'")),
						},
						{
							Expression:  *ast.NewConstExpr(ast.NewString("'a'")),
							SubRangeEnd: ast.NewConstExpr(ast.NewString("'z'")),
						},
					},
				},
			)},
			// AlphaNum = Alpha + Numeric;
			{Ident: ast.Ident("AlphaNum"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm("Alpha"),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: *ast.NewTerm("Numeric")},
					},
				},
			)},
		},
	)
}
