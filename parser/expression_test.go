package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Expression) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseExpression()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"variable",
		[]rune(`X`),
		asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("X")}),
	)

	run(
		"true constant of integer #1",
		[]rune(`7`),
		&ast.Expression{
			SimpleExpression: &ast.SimpleExpression{
				Term: &ast.Term{
					Factor: &ast.Number{
						ValueFactor: ast.ValueFactor{Value: "7"},
					},
				},
			},
		},
	)
	run(
		"true constant of integer #2",
		[]rune(`7`),
		asttest.NewExpression(asttest.NewNumber("7")),
	)

	run(
		"true constant of string",
		[]rune(`'abc'`),
		&ast.Expression{
			SimpleExpression: &ast.SimpleExpression{
				Term: &ast.Term{
					Factor: &ast.String{
						ValueFactor: ast.ValueFactor{Value: "'abc'"},
					},
				},
			},
		},
	)

	run(
		"address of variable",
		[]rune(`@X`),
		asttest.NewExpression(
			&ast.Address{
				Designator: &ast.Designator{
					QualId: &ast.QualId{Ident: asttest.NewIdent("X")},
				},
			},
		),
	)

	run(
		"integer constant",
		[]rune(`15`),
		asttest.NewExpression(&ast.Number{ValueFactor: ast.ValueFactor{Value: "15"}}),
	)

	run(
		"variable#2",
		[]rune(`InterestRate`),
		asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("InterestRate")}),
	)

	run(
		"function call",
		[]rune(`Calc(X,Y)`),
		asttest.NewExpression(
			&ast.DesignatorFactor{
				Designator: &ast.Designator{
					QualId: &ast.QualId{Ident: asttest.NewIdent("Calc")},
				},
				ExprList: ast.ExprList{
					asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("X")}),
					asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("Y")}),
				},
			},
		),
	)

	run(
		"quotient of Z and ( 1 - Z )",
		[]rune(`Z / (1 - Z)`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("Z"),
				MulOpFactors: []*ast.MulOpFactor{
					{
						MulOp: "/",
						Factor: &ast.Parentheses{
							Expression: &ast.Expression{
								SimpleExpression: &ast.SimpleExpression{
									Term: &ast.Term{Factor: asttest.NewNumber("1")},
									AddOpTerms: []*ast.AddOpTerm{
										{AddOp: "-", Term: asttest.NewTerm("Z")},
									},
								},
							},
						},
					},
				},
			},
		),
	)

	run(
		"Boolean #1",
		[]rune(`X = 1.5`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("1.5")),
				},
			},
		},
	)

	run(
		"Boolean #2",
		[]rune(`C in Range1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("C"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "IN",
					SimpleExpression: asttest.NewSimpleExpression("Range1"),
				},
			},
		},
	)

	run(
		"negation of a Boolean",
		[]rune(`not Done`),
		asttest.NewExpression(
			&ast.Not{
				Factor: asttest.NewDesignatorFactor("Done"),
			},
		),
	)

	run(
		"set",
		[]rune(`['a','b','c']`),
		asttest.NewExpression(
			&ast.SetConstructor{
				SetElements: []*ast.SetElement{
					asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'a'"))),
					asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'b'"))),
					asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'c'"))),
				},
			},
		),
	)

	// run(
	// 	"value typecast",
	// 	[]rune(`Char(48)`),
	// 	asttest.NewExpression(
	// 		&ast.TypeCast{
	// 			TypeId:     &ast.TypeId{Ident: "Char"},
	// 			Expression: *asttest.NewExpression(asttest.NewNumber("48")),
	// 		},
	// 	),
	// )
	run(
		"value typecast",
		[]rune(`Char(48)`),
		asttest.NewExpression(
			&ast.DesignatorFactor{
				Designator: &ast.Designator{
					QualId: &ast.QualId{Ident: asttest.NewIdent("Char")},
				},
				ExprList: ast.ExprList{
					asttest.NewExpression(asttest.NewNumber("48")),
				},
			},
		),
	)

	run(
		"Binary arithmetic operators +",
		[]rune(`X + Y`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators -",
		[]rune(`Result - 1`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("Result")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: asttest.NewTerm(asttest.NewNumber("1"))},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators *",
		[]rune(`P * InterestRate`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("P"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: asttest.NewDesignatorFactor("InterestRate")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators /",
		[]rune(`X / 2`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "/", Factor: asttest.NewNumber("2")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators div",
		[]rune(`Total div UnitSize`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("Total"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "DIV", Factor: asttest.NewDesignatorFactor("UnitSize")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators mod",
		[]rune(`Y mod 6`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("Y"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "MOD", Factor: asttest.NewNumber("6")},
				},
			},
		),
	)

	run(
		"Unary arithmetic operators +",
		[]rune(`+7`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("+"),
				Term:    asttest.NewTerm(asttest.NewNumber("7")),
			},
		),
	)

	run(
		"Unary arithmetic operators -",
		[]rune(`-X`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("-"),
				Term:    asttest.NewTerm("X"),
			},
		),
	)

	// Boolean operators

	run(
		"Boolean operators not",
		[]rune(`not (C in MySet)`),
		asttest.NewExpression(
			&ast.Not{
				Factor: &ast.Parentheses{
					Expression: &ast.Expression{
						SimpleExpression: asttest.NewSimpleExpression("C"),
						RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
							{
								RelOp:            "IN",
								SimpleExpression: asttest.NewSimpleExpression("MySet"),
							},
						},
					},
				},
			},
		),
	)

	run(
		"Boolean operators and",
		[]rune(`Done and (Total > 0)`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("Done"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "AND", Factor: &ast.Parentheses{
						Expression: &ast.Expression{
							SimpleExpression: asttest.NewSimpleExpression("Total"),
							RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
								{
									RelOp:            ">",
									SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0")),
								},
							},
						},
					}},
				},
			},
		),
	)

	run(
		"Boolean operators or",
		[]rune(`A or B`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("A")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "OR", Term: asttest.NewTerm("B")},
				},
			},
		),
	)

	run(
		"Boolean operators xor",
		[]rune(`A xor B`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("A")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "XOR", Term: asttest.NewTerm("B")},
				},
			},
		),
	)

	// Logical (bitwise) operators

	run(
		"Logical (bitwise) operators not",
		[]rune(`not X`),
		asttest.NewExpression(
			&ast.Not{
				Factor: asttest.NewDesignatorFactor("X"),
			},
		),
	)

	run(
		"Logical (bitwise) operators and",
		[]rune(`X and Y`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "AND", Factor: asttest.NewDesignatorFactor("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators or",
		[]rune(`X or Y`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "OR", Term: asttest.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators xor",
		[]rune(`X xor Y`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "XOR", Term: asttest.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators shl",
		[]rune(`X shl 2`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "SHL", Factor: asttest.NewNumber("2")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators shr",
		[]rune(`X shr I`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "SHR", Factor: asttest.NewDesignatorFactor("I")},
				},
			},
		),
	)

	// String operators

	run(
		"String operators +",
		[]rune(`S + '. '`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("S")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm(asttest.NewString("'. '"))},
				},
			},
		),
	)

	// Character-pointer operators

	run(
		"Character-pointer operators +",
		[]rune(`P + I`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("P")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm("I")},
				},
			},
		),
	)

	run(
		"Character-pointer operators -",
		[]rune(`P - Q`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("P")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: asttest.NewTerm("Q")},
				},
			},
		),
	)

	run(
		"Character-pointer operators ^",
		[]rune(`P^`),
		asttest.NewExpression(
			&ast.Designator{
				QualId: &ast.QualId{Ident: asttest.NewIdent("P")},
				Items: []ast.DesignatorItem{
					&ast.DesignatorItemDereference{},
				},
			},
		),
	)

	run(
		"Character-pointer operators =",
		[]rune(`P = Q`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("P"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: asttest.NewSimpleExpression("Q"),
				},
			},
		},
	)

	run(
		"Character-pointer operators <>",
		[]rune(`P <> Q`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("P"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: asttest.NewSimpleExpression("Q"),
				},
			},
		},
	)

	// Set operators

	run(
		"Set operators +",
		[]rune(`Set1 + Set2`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("Set1")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm("Set2")},
				},
			},
		),
	)

	run(
		"Set operators -",
		[]rune(`S - T`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor("S")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: asttest.NewTerm("T")},
				},
			},
		),
	)

	run(
		"Set operators *",
		[]rune(`S * T`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor("S"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: asttest.NewDesignatorFactor("T")},
				},
			},
		),
	)

	run(
		"Set operators <=",
		[]rune(`Q <= MySet`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("Q"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<=",
					SimpleExpression: asttest.NewSimpleExpression("MySet"),
				},
			},
		},
	)

	run(
		"Set operators >=",
		[]rune(`S1 >= S2`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("S1"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">=",
					SimpleExpression: asttest.NewSimpleExpression("S2"),
				},
			},
		},
	)

	run(
		"Set operators =",
		[]rune(`S2 = MySet`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("S2"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: asttest.NewSimpleExpression("MySet"),
				},
			},
		},
	)

	run(
		"Set operators <>",
		[]rune(`MySet <> S1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("MySet"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: asttest.NewSimpleExpression("S1"),
				},
			},
		},
	)

	run(
		"Set operators in",
		[]rune(`A in Set1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("A"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "IN",
					SimpleExpression: asttest.NewSimpleExpression("Set1"),
				},
			},
		},
	)

	// Relational operators

	run(
		"Relational operators =",
		[]rune(`I = Max`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("I"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: asttest.NewSimpleExpression("Max"),
				},
			},
		},
	)

	run(
		"Relational operators <>",
		[]rune(`X <> Y`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: asttest.NewSimpleExpression("Y"),
				},
			},
		},
	)

	run(
		"Relational operators <",
		[]rune(`X < Y`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<",
					SimpleExpression: asttest.NewSimpleExpression("Y"),
				},
			},
		},
	)

	run(
		"Relational operators >",
		[]rune(`Len > 0`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("Len"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0")),
				},
			},
		},
	)
	run(
		"Relational operators <=",
		[]rune(`Cnt <= I`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("Cnt"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<=",
					SimpleExpression: asttest.NewSimpleExpression("I"),
				},
			},
		},
	)

	run(
		"Relational operators >=",
		[]rune(`I >= 1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression("I"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">=",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("1")),
				},
			},
		},
	)
}

func TestSimpleExpression(t *testing.T) {
	run := func(name string, text []rune, expected *ast.SimpleExpression) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseSimpleExpression()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"Set operators -",
		[]rune(`S - T`),
		&ast.SimpleExpression{
			Term: &ast.Term{Factor: asttest.NewDesignatorFactor("S")},
			AddOpTerms: []*ast.AddOpTerm{
				{AddOp: "-", Term: asttest.NewTerm("T")},
			},
		},
	)
}

func TestSetConstructor(t *testing.T) {
	run := func(name string, text []rune, expected *ast.SetConstructor) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseSetConstructor()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"set",
		[]rune(`['a','b','c']`),
		&ast.SetConstructor{
			SetElements: []*ast.SetElement{
				asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'a'"))),
				asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'b'"))),
				asttest.NewSetElement(asttest.NewExpression(asttest.NewString("'c'"))),
			},
		},
	)
}

func TestFactor(t *testing.T) {
	run := func(name string, text []rune, expected ast.Factor) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseFactor()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"function call",
		[]rune(`Calc(X,Y)`),
		&ast.DesignatorFactor{
			Designator: &ast.Designator{
				QualId: &ast.QualId{Ident: asttest.NewIdent("Calc")},
			},
			ExprList: ast.ExprList{
				asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("X")}),
				asttest.NewExpression(&ast.QualId{Ident: asttest.NewIdent("Y")}),
			},
		},
	)
}
