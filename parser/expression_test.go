package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
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
		ast.NewExpression(&ast.QualId{Ident: ast.Ident("X")}),
	)

	run(
		"address of variable",
		[]rune(`@X`),
		ast.NewExpression(
			&ast.Address{
				Designator: ast.Designator{
					QualId: ast.QualId{Ident: ast.Ident("X")},
				},
			},
		),
	)

	run(
		"integer constant",
		[]rune(`15`),
		ast.NewExpression(&ast.Number{ValueFactor: ast.ValueFactor{Value: "15"}}),
	)

	run(
		"variable#2",
		[]rune(`InterestRate`),
		ast.NewExpression(&ast.QualId{Ident: ast.Ident("InterestRate")}),
	)

	run(
		"function call",
		[]rune(`Calc(X,Y)`),
		ast.NewExpression(
			&ast.DesignatorFactor{
				Designator: ast.Designator{
					QualId: ast.QualId{Ident: ast.Ident("Calc")},
				},
				ExprList: ast.ExprList{
					ast.NewExpression(&ast.QualId{Ident: ast.Ident("X")}),
					ast.NewExpression(&ast.QualId{Ident: ast.Ident("Y")}),
				},
			},
		),
	)

	run(
		"quotient of Z and ( 1 - Z )",
		[]rune(`Z / (1 - Z)`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("Z"),
				MulOpFactors: []*ast.MulOpFactor{
					{
						MulOp: "/",
						Factor: &ast.Parentheses{
							Expression: ast.Expression{
								SimpleExpression: ast.SimpleExpression{
									Term: ast.Term{Factor: ast.NewNumber("1")},
									AddOpTerms: []*ast.AddOpTerm{
										{AddOp: "-", Term: *ast.NewTerm("Z")},
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
			SimpleExpression: *ast.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: *ast.NewSimpleExpression(ast.NewNumber("1.5")),
				},
			},
		},
	)

	run(
		"Boolean #2",
		[]rune(`C in Range1`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("C"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "IN",
					SimpleExpression: *ast.NewSimpleExpression("Range1"),
				},
			},
		},
	)

	run(
		"negation of a Boolean",
		[]rune(`not Done`),
		ast.NewExpression(
			&ast.Not{
				Factor: ast.NewDesignatorFactor("Done"),
			},
		),
	)

	run(
		"set",
		[]rune(`['a','b','c']`),
		ast.NewExpression(
			&ast.SetConstructor{
				SetElements: []*ast.SetElement{
					ast.NewSetElement(ast.NewExpression(ast.NewString("'a'"))),
					ast.NewSetElement(ast.NewExpression(ast.NewString("'b'"))),
					ast.NewSetElement(ast.NewExpression(ast.NewString("'c'"))),
				},
			},
		),
	)

	// run(
	// 	"value typecast",
	// 	[]rune(`Char(48)`),
	// 	ast.NewExpression(
	// 		&ast.TypeCast{
	// 			TypeId:     &ast.TypeId{Ident: "Char"},
	// 			Expression: *ast.NewExpression(ast.NewNumber("48")),
	// 		},
	// 	),
	// )
	run(
		"value typecast",
		[]rune(`Char(48)`),
		ast.NewExpression(
			&ast.DesignatorFactor{
				Designator: ast.Designator{
					QualId: ast.QualId{Ident: ast.Ident("Char")},
				},
				ExprList: ast.ExprList{
					ast.NewExpression(ast.NewNumber("48")),
				},
			},
		),
	)

	run(
		"Binary arithmetic operators +",
		[]rune(`X + Y`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: *ast.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators -",
		[]rune(`Result - 1`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("Result")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: *ast.NewTerm(ast.NewNumber("1"))},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators *",
		[]rune(`P * InterestRate`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("P"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: ast.NewDesignatorFactor("InterestRate")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators /",
		[]rune(`X / 2`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "/", Factor: ast.NewNumber("2")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators div",
		[]rune(`Total div UnitSize`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("Total"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "DIV", Factor: ast.NewDesignatorFactor("UnitSize")},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators mod",
		[]rune(`Y mod 6`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("Y"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "MOD", Factor: ast.NewNumber("6")},
				},
			},
		),
	)

	run(
		"Unary arithmetic operators +",
		[]rune(`+7`),
		ast.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("+"),
				Term:    *ast.NewTerm(ast.NewNumber("7")),
			},
		),
	)

	run(
		"Unary arithmetic operators -",
		[]rune(`-X`),
		ast.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("-"),
				Term:    *ast.NewTerm("X"),
			},
		),
	)

	// Boolean operators

	run(
		"Boolean operators not",
		[]rune(`not (C in MySet)`),
		ast.NewExpression(
			&ast.Not{
				Factor: &ast.Parentheses{
					Expression: ast.Expression{
						SimpleExpression: *ast.NewSimpleExpression("C"),
						RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
							{
								RelOp:            "IN",
								SimpleExpression: *ast.NewSimpleExpression("MySet"),
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
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("Done"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "AND", Factor: &ast.Parentheses{
						Expression: ast.Expression{
							SimpleExpression: *ast.NewSimpleExpression("Total"),
							RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
								{
									RelOp:            ">",
									SimpleExpression: *ast.NewSimpleExpression(ast.NewNumber("0")),
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
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("A")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "OR", Term: *ast.NewTerm("B")},
				},
			},
		),
	)

	run(
		"Boolean operators xor",
		[]rune(`A xor B`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("A")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "XOR", Term: *ast.NewTerm("B")},
				},
			},
		),
	)

	// Logical (bitwise) operators

	run(
		"Logical (bitwise) operators not",
		[]rune(`not X`),
		ast.NewExpression(
			&ast.Not{
				Factor: ast.NewDesignatorFactor("X"),
			},
		),
	)

	run(
		"Logical (bitwise) operators and",
		[]rune(`X and Y`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "AND", Factor: ast.NewDesignatorFactor("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators or",
		[]rune(`X or Y`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "OR", Term: *ast.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators xor",
		[]rune(`X xor Y`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("X")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "XOR", Term: *ast.NewTerm("Y")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators shl",
		[]rune(`X shl 2`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "SHL", Factor: ast.NewNumber("2")},
				},
			},
		),
	)

	run(
		"Logical (bitwise) operators shr",
		[]rune(`X shr I`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("X"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "SHR", Factor: ast.NewDesignatorFactor("I")},
				},
			},
		),
	)

	// String operators

	run(
		"String operators +",
		[]rune(`S + '. '`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("S")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: *ast.NewTerm(ast.NewString("'. '"))},
				},
			},
		),
	)

	// Character-pointer operators

	run(
		"Character-pointer operators +",
		[]rune(`P + I`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("P")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: *ast.NewTerm("I")},
				},
			},
		),
	)

	run(
		"Character-pointer operators -",
		[]rune(`P - Q`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("P")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: *ast.NewTerm("Q")},
				},
			},
		),
	)

	run(
		"Character-pointer operators ^",
		[]rune(`P^`),
		ast.NewExpression(
			&ast.Designator{
				QualId: ast.QualId{Ident: ast.Ident("P")},
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
			SimpleExpression: *ast.NewSimpleExpression("P"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: *ast.NewSimpleExpression("Q"),
				},
			},
		},
	)

	run(
		"Character-pointer operators <>",
		[]rune(`P <> Q`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("P"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: *ast.NewSimpleExpression("Q"),
				},
			},
		},
	)

	// Set operators

	run(
		"Set operators +",
		[]rune(`Set1 + Set2`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("Set1")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: *ast.NewTerm("Set2")},
				},
			},
		),
	)

	run(
		"Set operators -",
		[]rune(`S - T`),
		ast.NewExpression(
			&ast.SimpleExpression{
				Term: ast.Term{Factor: ast.NewDesignatorFactor("S")},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: *ast.NewTerm("T")},
				},
			},
		),
	)

	run(
		"Set operators *",
		[]rune(`S * T`),
		ast.NewExpression(
			&ast.Term{
				Factor: ast.NewDesignatorFactor("S"),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: ast.NewDesignatorFactor("T")},
				},
			},
		),
	)

	run(
		"Set operators <=",
		[]rune(`Q <= MySet`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("Q"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<=",
					SimpleExpression: *ast.NewSimpleExpression("MySet"),
				},
			},
		},
	)

	run(
		"Set operators >=",
		[]rune(`S1 >= S2`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("S1"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">=",
					SimpleExpression: *ast.NewSimpleExpression("S2"),
				},
			},
		},
	)

	run(
		"Set operators =",
		[]rune(`S2 = MySet`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("S2"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: *ast.NewSimpleExpression("MySet"),
				},
			},
		},
	)

	run(
		"Set operators <>",
		[]rune(`MySet <> S1`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("MySet"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: *ast.NewSimpleExpression("S1"),
				},
			},
		},
	)

	run(
		"Set operators in",
		[]rune(`A in Set1`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("A"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "IN",
					SimpleExpression: *ast.NewSimpleExpression("Set1"),
				},
			},
		},
	)

	// Relational operators

	run(
		"Relational operators =",
		[]rune(`I = Max`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("I"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: *ast.NewSimpleExpression("Max"),
				},
			},
		},
	)

	run(
		"Relational operators <>",
		[]rune(`X <> Y`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: *ast.NewSimpleExpression("Y"),
				},
			},
		},
	)

	run(
		"Relational operators <",
		[]rune(`X < Y`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("X"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<",
					SimpleExpression: *ast.NewSimpleExpression("Y"),
				},
			},
		},
	)

	run(
		"Relational operators >",
		[]rune(`Len > 0`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("Len"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">",
					SimpleExpression: *ast.NewSimpleExpression(ast.NewNumber("0")),
				},
			},
		},
	)
	run(
		"Relational operators <=",
		[]rune(`Cnt <= I`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("Cnt"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<=",
					SimpleExpression: *ast.NewSimpleExpression("I"),
				},
			},
		},
	)

	run(
		"Relational operators >=",
		[]rune(`I >= 1`),
		&ast.Expression{
			SimpleExpression: *ast.NewSimpleExpression("I"),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">=",
					SimpleExpression: *ast.NewSimpleExpression(ast.NewNumber("1")),
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
			Term: ast.Term{Factor: ast.NewDesignatorFactor("S")},
			AddOpTerms: []*ast.AddOpTerm{
				{AddOp: "-", Term: *ast.NewTerm("T")},
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
				ast.NewSetElement(ast.NewExpression(ast.NewString("'a'"))),
				ast.NewSetElement(ast.NewExpression(ast.NewString("'b'"))),
				ast.NewSetElement(ast.NewExpression(ast.NewString("'c'"))),
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
			Designator: ast.Designator{
				QualId: ast.QualId{Ident: ast.Ident("Calc")},
			},
			ExprList: ast.ExprList{
				ast.NewExpression(&ast.QualId{Ident: ast.Ident("X")}),
				ast.NewExpression(&ast.QualId{Ident: ast.Ident("Y")}),
			},
		},
	)
}
