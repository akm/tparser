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
			Ident: asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ConstSection{
						{Ident: asttest.NewIdent("MaxValue"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("237"))},
						{Ident: asttest.NewIdent("Message1"), ConstExpr: asttest.NewConstExpr(asttest.NewString("'Out of memory'"))},
						{Ident: asttest.NewIdent("Max"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("100")), Type: asttest.NewOrdIdent("Integer")},
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
			{Ident: asttest.NewIdent("MaxValue"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("237"))},
		},
	)
	run(
		"number const",
		[]rune(`CONST Max: Integer = 100;`),
		ast.ConstSection{
			{Ident: asttest.NewIdent("Max"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("100")), Type: asttest.NewOrdIdent("Integer")},
		},
	)
	run(
		"message as identifier",
		[]rune(`CONST Message = 'Out of memory';`),
		ast.ConstSection{
			{Ident: asttest.NewIdent("Message"), ConstExpr: asttest.NewConstExpr(asttest.NewString("'Out of memory'"))},
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
			{Ident: asttest.NewIdent("Min"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("0"))},
			{Ident: asttest.NewIdent("Max"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("100"))},
			// Center = (Max - Min) div 2;
			{Ident: asttest.NewIdent("Center"), ConstExpr: asttest.NewConstExpr(
				&ast.Term{
					Factor: &ast.Parentheses{
						Expression: &ast.Expression{
							SimpleExpression: &ast.SimpleExpression{
								Term: asttest.NewTerm("Max"),
								AddOpTerms: []*ast.AddOpTerm{
									{AddOp: "-", Term: asttest.NewTerm("Min")},
								},
							},
						},
					},
					MulOpFactors: []*ast.MulOpFactor{
						{MulOp: "DIV", Factor: asttest.NewNumber("2")},
					},
				},
			)},
			// Beta = Chr(225);
			{Ident: asttest.NewIdent("Beta"), ConstExpr: asttest.NewConstExpr(
				&ast.DesignatorFactor{
					Designator: &ast.Designator{
						QualId: &ast.QualId{Ident: asttest.NewIdent("Chr")},
					},
					ExprList: ast.ExprList{
						asttest.NewConstExpr(asttest.NewNumber("225")),
					},
				},
			)},
			// NumChars = Ord('Z') - Ord('A') + 1;
			{Ident: asttest.NewIdent("NumChars"), ConstExpr: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(
						&ast.DesignatorFactor{
							Designator: &ast.Designator{
								QualId: &ast.QualId{Ident: asttest.NewIdent("Ord")},
							},
							ExprList: ast.ExprList{
								asttest.NewConstExpr(asttest.NewString("'Z'")),
							},
						},
					),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "-",
							Term: asttest.NewTerm(
								&ast.DesignatorFactor{
									Designator: &ast.Designator{
										QualId: &ast.QualId{Ident: asttest.NewIdent("Ord")},
									},
									ExprList: ast.ExprList{
										asttest.NewConstExpr(asttest.NewString("'A'")),
									},
								},
							),
						},
						{
							AddOp: "+",
							Term:  asttest.NewTerm(asttest.NewNumber("1")),
						},
					},
				},
			)},
			// Message = 'Out of memory';
			{Ident: asttest.NewIdent("Message"), ConstExpr: asttest.NewConstExpr(asttest.NewString("'Out of memory'"))},
			// ErrStr = ' Error: ' + Message + '. ';
			{Ident: asttest.NewIdent("ErrStr"), ConstExpr: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(asttest.NewString("' Error: '")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "+",
							Term:  asttest.NewTerm("Message"),
						},
						{
							AddOp: "+",
							Term:  asttest.NewTerm(asttest.NewString("'. '")),
						},
					},
				},
			)},
			// ErrPos = 80 - Length(ErrStr) div 2;
			{Ident: asttest.NewIdent("ErrPos"), ConstExpr: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(asttest.NewNumber("80")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "-",
							Term: &ast.Term{
								Factor: &ast.DesignatorFactor{
									Designator: &ast.Designator{
										QualId: &ast.QualId{Ident: asttest.NewIdent("Length")},
									},
									ExprList: ast.ExprList{
										asttest.NewConstExpr("ErrStr"),
									},
								},
								MulOpFactors: []*ast.MulOpFactor{
									{MulOp: "DIV", Factor: asttest.NewNumber("2")},
								},
							},
						},
					},
				},
			)},
			// Ln10 = 2.302585092994045684;
			{Ident: asttest.NewIdent("Ln10"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("2.302585092994045684"))},
			// Ln10R = 1 / Ln10;
			{Ident: asttest.NewIdent("Ln10R"), ConstExpr: asttest.NewConstExpr(
				&ast.Term{
					Factor: asttest.NewNumber("1"),
					MulOpFactors: []*ast.MulOpFactor{
						{MulOp: "/", Factor: asttest.NewDesignatorFactor("Ln10")},
					},
				},
			)},
			// Numeric = ['0'..'9'];
			{Ident: asttest.NewIdent("Numeric"), ConstExpr: asttest.NewConstExpr(
				&ast.SetConstructor{
					SetElements: []*ast.SetElement{
						{
							Expression:  asttest.NewConstExpr(asttest.NewString("'0'")),
							SubRangeEnd: asttest.NewConstExpr(asttest.NewString("'9'")),
						},
					},
				},
			)},
			// Alpha = ['A'..'Z', 'a'..'z'];
			{Ident: asttest.NewIdent("Alpha"), ConstExpr: asttest.NewConstExpr(
				&ast.SetConstructor{
					SetElements: []*ast.SetElement{
						{
							Expression:  asttest.NewConstExpr(asttest.NewString("'A'")),
							SubRangeEnd: asttest.NewConstExpr(asttest.NewString("'Z'")),
						},
						{
							Expression:  asttest.NewConstExpr(asttest.NewString("'a'")),
							SubRangeEnd: asttest.NewConstExpr(asttest.NewString("'z'")),
						},
					},
				},
			)},
			// AlphaNum = Alpha + Numeric;
			{Ident: asttest.NewIdent("AlphaNum"), ConstExpr: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm("Alpha"),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: asttest.NewTerm("Numeric")},
					},
				},
			)},
		},
	)
}
