package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithConstSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
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
			Ident: ast.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ConstSection{
						{Ident: ast.NewIdent("MaxValue"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("237"))},
						{Ident: ast.NewIdent("Message1"), ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'"))},
						{Ident: ast.NewIdent("Max"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("100")), Type: ast.NewOrdIdent("Integer")},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestConstSectionl(t *testing.T) {
	run := func(name string, text []rune, expected ast.ConstSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseConstSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"number const",
		[]rune(`CONST MaxValue = 237;`),
		ast.ConstSection{
			{Ident: ast.NewIdent("MaxValue"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("237"))},
		},
	)
	run(
		"number const",
		[]rune(`CONST Max: Integer = 100;`),
		ast.ConstSection{
			{Ident: ast.NewIdent("Max"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("100")), Type: ast.NewOrdIdent("Integer")},
		},
	)
	run(
		"message as identifier",
		[]rune(`CONST Message = 'Out of memory';`),
		ast.ConstSection{
			{Ident: ast.NewIdent("Message"), ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'"))},
		},
	)

	run(
		"examples in Language Guide",
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
			{Ident: ast.NewIdent("Min"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("0"))},
			{Ident: ast.NewIdent("Max"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("100"))},
			// Center = (Max - Min) div 2;
			{Ident: ast.NewIdent("Center"), ConstExpr: *ast.NewConstExpr(
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
			{Ident: ast.NewIdent("Beta"), ConstExpr: *ast.NewConstExpr(
				&ast.DesignatorFactor{
					Designator: ast.Designator{
						QualId: ast.QualId{Ident: ast.NewIdent("Chr")},
					},
					ExprList: ast.ExprList{
						ast.NewConstExpr(ast.NewNumber("225")),
					},
				},
			)},
			// NumChars = Ord('Z') - Ord('A') + 1;
			{Ident: ast.NewIdent("NumChars"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm(
						&ast.DesignatorFactor{
							Designator: ast.Designator{
								QualId: ast.QualId{Ident: ast.NewIdent("Ord")},
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
										QualId: ast.QualId{Ident: ast.NewIdent("Ord")},
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
			{Ident: ast.NewIdent("Message"), ConstExpr: *ast.NewConstExpr(ast.NewString("'Out of memory'"))},
			// ErrStr = ' Error: ' + Message + '. ';
			{Ident: ast.NewIdent("ErrStr"), ConstExpr: *ast.NewConstExpr(
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
			{Ident: ast.NewIdent("ErrPos"), ConstExpr: *ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm(ast.NewNumber("80")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "-",
							Term: ast.Term{
								Factor: &ast.DesignatorFactor{
									Designator: ast.Designator{
										QualId: ast.QualId{Ident: ast.NewIdent("Length")},
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
			{Ident: ast.NewIdent("Ln10"), ConstExpr: *ast.NewConstExpr(ast.NewNumber("2.302585092994045684"))},
			// Ln10R = 1 / Ln10;
			{Ident: ast.NewIdent("Ln10R"), ConstExpr: *ast.NewConstExpr(
				&ast.Term{
					Factor: ast.NewNumber("1"),
					MulOpFactors: []*ast.MulOpFactor{
						{MulOp: "/", Factor: ast.NewDesignatorFactor("Ln10")},
					},
				},
			)},
			// Numeric = ['0'..'9'];
			{Ident: ast.NewIdent("Numeric"), ConstExpr: *ast.NewConstExpr(
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
			{Ident: ast.NewIdent("Alpha"), ConstExpr: *ast.NewConstExpr(
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
			{Ident: ast.NewIdent("AlphaNum"), ConstExpr: *ast.NewConstExpr(
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
