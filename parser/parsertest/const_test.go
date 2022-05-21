package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithConstSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				asttest.ClearUnitDeclarationMap(res)
				asttest.ClearLocations(t, res)
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
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseConstSection(true)
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				if !assert.Equal(t, expected, res) {
					for idx, decl := range res {
						assert.Equal(t, expected[idx].Ident.Name, decl.Ident.Name)
						if !assert.Equal(t, expected[idx], decl, "declaration %d - %s", idx, decl.Ident.Name) {
							assert.Equal(t, expected[idx].Ident, decl.Ident)
							assert.Equal(t, expected[idx].Type, decl.Type)
							assert.Equal(t, expected[idx].ConstExpr, decl.ConstExpr)
							// assert.Equal(t, expected[idx].ConstExpr.String(), decl.ConstExpr.String())
							// spew.Dump(expected[idx].ConstExpr)
							// spew.Dump(decl.ConstExpr)
						}
					}
				}
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
		func() ast.ConstSection {
			min := &ast.ConstantDecl{Ident: asttest.NewIdent("Min"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("0"))}
			max := &ast.ConstantDecl{Ident: asttest.NewIdent("Max"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("100"))}

			// Message = 'Out of memory';
			message := &ast.ConstantDecl{Ident: asttest.NewIdent("Message"), ConstExpr: asttest.NewConstExpr(asttest.NewString("'Out of memory'"))}
			// ErrStr = ' Error: ' + Message + '. ';
			errStr := &ast.ConstantDecl{Ident: asttest.NewIdent("ErrStr"), ConstExpr: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(asttest.NewString("' Error: '")),
					AddOpTerms: []*ast.AddOpTerm{
						{
							AddOp: "+",
							Term:  asttest.NewTerm(asttest.NewQualId("Message", message.ToDeclarations()[0])),
						},
						{
							AddOp: "+",
							Term:  asttest.NewTerm(asttest.NewString("'. '")),
						},
					},
				},
			)}

			// Ln10 = 2.302585092994045684;
			ln10 := &ast.ConstantDecl{Ident: asttest.NewIdent("Ln10"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("2.302585092994045684"))}

			// Numeric = ['0'..'9'];
			numeric := &ast.ConstantDecl{Ident: asttest.NewIdent("Numeric"), ConstExpr: asttest.NewConstExpr(
				&ast.SetConstructor{
					SetElements: []*ast.SetElement{
						{
							Expression:  asttest.NewConstExpr(asttest.NewString("'0'")),
							SubRangeEnd: asttest.NewConstExpr(asttest.NewString("'9'")),
						},
					},
				},
			)}
			// Alpha = ['A'..'Z', 'a'..'z'];
			alpha := &ast.ConstantDecl{Ident: asttest.NewIdent("Alpha"), ConstExpr: asttest.NewConstExpr(
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
			)}

			return ast.ConstSection{
				min,
				max,
				// Center = (Max - Min) div 2;
				{Ident: asttest.NewIdent("Center"), ConstExpr: asttest.NewConstExpr(
					&ast.Term{
						Factor: &ast.Parentheses{
							Expression: &ast.Expression{
								SimpleExpression: &ast.SimpleExpression{
									Term: asttest.NewTerm(asttest.NewQualId("Max", max.ToDeclarations()[0])),
									AddOpTerms: []*ast.AddOpTerm{
										{AddOp: "-", Term: asttest.NewTerm(asttest.NewQualId("Min", min.ToDeclarations()[0]))},
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
						Designator: &ast.Designator{QualId: asttest.NewQualId("Chr")},
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
								Designator: &ast.Designator{QualId: asttest.NewQualId("Ord")},
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
										Designator: &ast.Designator{QualId: asttest.NewQualId("Ord")},
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
				message,
				errStr,
				// ErrPos = 80 - Length(ErrStr) div 2;
				{Ident: asttest.NewIdent("ErrPos"), ConstExpr: asttest.NewConstExpr(
					&ast.SimpleExpression{
						Term: asttest.NewTerm(asttest.NewNumber("80")),
						AddOpTerms: []*ast.AddOpTerm{
							{
								AddOp: "-",
								Term: &ast.Term{
									Factor: &ast.DesignatorFactor{
										Designator: &ast.Designator{QualId: asttest.NewQualId("Length")},
										ExprList: ast.ExprList{
											asttest.NewConstExpr(asttest.NewQualId("ErrStr", errStr.ToDeclarations()[0])),
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
				ln10,
				// Ln10R = 1 / Ln10;
				{Ident: asttest.NewIdent("Ln10R"), ConstExpr: asttest.NewConstExpr(
					&ast.Term{
						Factor: asttest.NewNumber("1"),
						MulOpFactors: []*ast.MulOpFactor{
							{MulOp: "/", Factor: asttest.NewDesignatorFactor(asttest.NewQualId("Ln10", ln10.ToDeclarations()[0]))},
						},
					},
				)},
				numeric,
				alpha,
				// AlphaNum = Alpha + Numeric;
				{Ident: asttest.NewIdent("AlphaNum"), ConstExpr: asttest.NewConstExpr(
					&ast.SimpleExpression{
						Term: asttest.NewTerm(asttest.NewQualId("Alpha", alpha.ToDeclarations()[0])),
						AddOpTerms: []*ast.AddOpTerm{
							{AddOp: "+", Term: asttest.NewTerm(asttest.NewQualId("Numeric", numeric.ToDeclarations()[0]))},
						},
					},
				)},
			}
		}(),
	)
}
