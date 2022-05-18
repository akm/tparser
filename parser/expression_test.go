package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/ext"
	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	run := func(name string, clearLocations bool, text []rune, expected *ast.Expression) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseExpression()
			if assert.NoError(t, err) {
				if clearLocations {
					asttest.ClearLocations(t, res)
				}
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"variable", false,
		[]rune(`X`),
		asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("X", asttest.NewIdentLocation(1, 1, 0, 1, 1, 1)))),
	)

	run(
		"true constant of integer #1", false,
		[]rune(`7`),
		&ast.Expression{
			SimpleExpression: &ast.SimpleExpression{
				Term: &ast.Term{
					Factor: &ast.NumberFactor{Value: "7"},
				},
			},
		},
	)
	run(
		"true constant of integer #2", false,
		[]rune(`7`),
		asttest.NewExpression(asttest.NewNumber("7")),
	)

	run(
		"true constant of string", false,
		[]rune(`'abc'`),
		&ast.Expression{
			SimpleExpression: &ast.SimpleExpression{
				Term: &ast.Term{
					Factor: &ast.StringFactor{Value: "'abc'"},
				},
			},
		},
	)

	run(
		"address of variable", false,
		[]rune(`@X`),
		asttest.NewExpression(
			&ast.Address{
				Designator: &ast.Designator{
					QualId: ast.NewQualId(nil, asttest.NewIdentRef("X", asttest.NewIdentLocation(1, 2, 1, 1, 2, 2))),
				},
			},
		),
	)

	run(
		"integer constant", false,
		[]rune(`15`),
		asttest.NewExpression(&ast.NumberFactor{Value: "15"}),
	)

	run(
		"variable#2", false,
		[]rune(`InterestRate`),
		asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("InterestRate", asttest.NewIdentLocation(1, 1, 0, 1, 12, 12)))),
	)

	run(
		"function call", false,
		[]rune(`Calc(X,Y)`),
		asttest.NewExpression(
			&ast.DesignatorFactor{
				Designator: &ast.Designator{
					QualId: ast.NewQualId(nil, asttest.NewIdentRef("Calc", asttest.NewIdentLocation(1, 1, 0, 1, 5, 4))),
				},
				ExprList: ast.ExprList{
					asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("X", asttest.NewIdentLocation(1, 6, 5, 7)))),
					asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("Y", asttest.NewIdentLocation(1, 8, 7, 9)))),
				},
			},
		),
	)

	run(
		"quotient of Z and ( 1 - Z )", false,
		[]rune(`Z / (1 - Z)`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Z", asttest.NewIdentLocation(1, 1, 0, 1, 2, 1))),
				MulOpFactors: []*ast.MulOpFactor{
					{
						MulOp: "/",
						Factor: &ast.Parentheses{
							Expression: &ast.Expression{
								SimpleExpression: &ast.SimpleExpression{
									Term: &ast.Term{Factor: asttest.NewNumber("1")},
									AddOpTerms: []*ast.AddOpTerm{
										{AddOp: "-", Term: asttest.NewTerm(asttest.NewIdent("Z", asttest.NewIdentLocation(1, 10, 9, 1, 11, 10)))},
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
		"Boolean #1", false,
		[]rune(`X = 1.5`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("X", asttest.NewIdentLocation(1, 1, 0, 1, 2, 1))),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "=",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("1.5")),
				},
			},
		},
	)

	run(
		"Boolean #2", false,
		[]rune(`C in Range1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("C", asttest.NewIdentLocation(1, 1, 0, 1, 2, 1))),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "IN",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("Range1", asttest.NewIdentLocation(1, 6, 5, 1, 11, 11))),
				},
			},
		},
	)

	run(
		"negation of a Boolean", false,
		[]rune(`not Done`),
		asttest.NewExpression(
			&ast.Not{
				Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Done", asttest.NewIdentLocation(1, 5, 4, 1, 8, 8))),
			},
		),
	)

	run(
		"set", false,
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
		"value typecast", false,
		[]rune(`Char(48)`),
		asttest.NewExpression(
			&ast.DesignatorFactor{
				Designator: &ast.Designator{
					QualId: ast.NewQualId(nil, asttest.NewIdentRef("Char", asttest.NewIdentLocation(1, 1, 0, 1, 5, 4))),
				},
				ExprList: ast.ExprList{
					asttest.NewExpression(asttest.NewNumber("48")),
				},
			},
		),
	)

	run(
		"Binary arithmetic operators +", false,
		[]rune(`X + Y`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("X", asttest.NewIdentLocation(1, 1, 0, 1, 2, 1)))},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm(asttest.NewIdent("Y", asttest.NewIdentLocation(1, 5, 4, 1, 5, 5)))},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators -", false,
		[]rune(`Result - 1`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Result", asttest.NewIdentLocation(1, 1, 0, 1, 7, 6)))},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "-", Term: asttest.NewTerm(asttest.NewNumber("1"))},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators *", false,
		[]rune(`P * InterestRate`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor(asttest.NewIdent("P", asttest.NewIdentLocation(1, 1, 0, 1, 2, 1))),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: asttest.NewDesignatorFactor(asttest.NewIdent("InterestRate", asttest.NewIdentLocation(1, 5, 4, 1, 16, 16)))},
				},
			},
		),
	)

	run(
		"Binary arithmetic operators /", true,
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
		"Binary arithmetic operators div", true,
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
		"Binary arithmetic operators mod", true,
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
		"Unary arithmetic operators +", false,
		[]rune(`+7`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("+"),
				Term:    asttest.NewTerm(asttest.NewNumber("7")),
			},
		),
	)

	run(
		"Unary arithmetic operators -", false,
		[]rune(`-X`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				UnaryOp: ext.StringPtr("-"),
				Term:    asttest.NewTerm(asttest.NewIdent("X", asttest.NewIdentLocation(1, 2, 1, 1, 2, 2))),
			},
		),
	)

	// Boolean operators

	run(
		"Boolean operators not", false,
		[]rune(`not (C in MySet)`),
		asttest.NewExpression(
			&ast.Not{
				Factor: &ast.Parentheses{
					Expression: &ast.Expression{
						SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("C", asttest.NewIdentLocation(1, 6, 5, 1, 7, 6))),
						RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
							{
								RelOp:            "IN",
								SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("MySet", asttest.NewIdentLocation(1, 11, 10, 1, 16, 15))),
							},
						},
					},
				},
			},
		),
	)

	run(
		"Boolean operators and", true,
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
		"Boolean operators or", true,
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
		"Boolean operators xor", true,
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
		"Logical (bitwise) operators not", true,
		[]rune(`not X`),
		asttest.NewExpression(
			&ast.Not{
				Factor: asttest.NewDesignatorFactor("X"),
			},
		),
	)

	run(
		"Logical (bitwise) operators and", true,
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
		"Logical (bitwise) operators or", true,
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
		"Logical (bitwise) operators xor", true,
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
		"Logical (bitwise) operators shl", true,
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
		"Logical (bitwise) operators shr", false,
		[]rune(`X shr I`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor(asttest.NewIdent("X", asttest.NewIdentLocation(1, 1, 0, 2, 1))),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "SHR", Factor: asttest.NewDesignatorFactor(asttest.NewIdent("I", asttest.NewIdentLocation(1, 7, 6, 1, 7, 7)))},
				},
			},
		),
	)

	// String operators

	run(
		"String operators +", true,
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
		"Character-pointer operators +", false,
		[]rune(`P + I`),
		asttest.NewExpression(
			&ast.SimpleExpression{
				Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("P", asttest.NewIdentLocation(1, 1, 0, 2, 1)))},
				AddOpTerms: []*ast.AddOpTerm{
					{AddOp: "+", Term: asttest.NewTerm(asttest.NewIdent("I", asttest.NewIdentLocation(1, 5, 4, 1, 5, 5)))},
				},
			},
		),
	)

	run(
		"Character-pointer operators -", true,
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
		"Character-pointer operators ^", false,
		[]rune(`P^`),
		asttest.NewExpression(
			&ast.Designator{
				QualId: ast.NewQualId(nil, asttest.NewIdentRef("P", asttest.NewIdentLocation(1, 1, 0, 2, 1))),
				Items: []ast.DesignatorItem{
					&ast.DesignatorItemDereference{},
				},
			},
		),
	)

	run(
		"Character-pointer operators =", true,
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
		"Character-pointer operators <>", false,
		[]rune(`P <> Q`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("P", asttest.NewIdentLocation(1, 1, 0, 2, 1))),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<>",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("Q", asttest.NewIdentLocation(1, 6, 5, 1, 6, 6))),
				},
			},
		},
	)

	// Set operators

	run(
		"Set operators +", true,
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
		"Set operators -", true,
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
		"Set operators *", false,
		[]rune(`S * T`),
		asttest.NewExpression(
			&ast.Term{
				Factor: asttest.NewDesignatorFactor(asttest.NewIdent("S", asttest.NewIdentLocation(1, 1, 0, 2, 1))),
				MulOpFactors: []*ast.MulOpFactor{
					{MulOp: "*", Factor: asttest.NewDesignatorFactor(asttest.NewIdent("T", asttest.NewIdentLocation(1, 5, 4, 1, 5, 5)))},
				},
			},
		),
	)

	run(
		"Set operators <=", true,
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
		"Set operators >=", true,
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
		"Set operators =", true,
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
		"Set operators <>", true,
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
		"Set operators in", true,
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
		"Relational operators =", true,
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
		"Relational operators <>", true,
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
		"Relational operators <", true,
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
		"Relational operators >", false,
		[]rune(`Len > 0`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("Len", asttest.NewIdentLocation(1, 1, 0, 1, 4, 3))),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            ">",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0")),
				},
			},
		},
	)
	run(
		"Relational operators <=", false,
		[]rune(`Cnt <= I`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("Cnt", asttest.NewIdentLocation(1, 1, 0, 1, 4, 3))),
			RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
				{
					RelOp:            "<=",
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("I", asttest.NewIdentLocation(1, 8, 7, 1, 8, 8))),
				},
			},
		},
	)

	run(
		"Relational operators >=", false,
		[]rune(`I >= 1`),
		&ast.Expression{
			SimpleExpression: asttest.NewSimpleExpression(asttest.NewIdent("I", asttest.NewIdentLocation(1, 1, 0, 2))),
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
			Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("S", asttest.NewIdentLocation(1, 1, 0, 2)))},
			AddOpTerms: []*ast.AddOpTerm{
				{AddOp: "-", Term: asttest.NewTerm(asttest.NewIdent("T", asttest.NewIdentLocation(1, 5, 4, 1, 5, 5)))},
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
				QualId: ast.NewQualId(nil, asttest.NewIdentRef("Calc", asttest.NewIdentLocation(1, 1, 0, 5))),
			},
			ExprList: ast.ExprList{
				asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("X", asttest.NewIdentLocation(1, 6, 5, 7)))),
				asttest.NewExpression(ast.NewQualId(nil, asttest.NewIdentRef("Y", asttest.NewIdentLocation(1, 8, 7, 9)))),
			},
		},
	)
}
