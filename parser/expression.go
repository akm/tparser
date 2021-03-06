package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseExprList(terminator token.Predicator) (ast.ExprList, error) {
	res := ast.ExprList{}
	for {
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res = append(res, expr)
		if p.CurrentToken().Is(terminator) {
			break
		}
		if _, err := p.Current(token.Symbol(',')); err != nil {
			return nil, err
		}
		p.NextToken()
	}
	return res, nil
}

func (p *Parser) ParseExpression() (*ast.Expression, error) {
	res := &ast.Expression{}
	se, err := p.ParseSimpleExpression()
	if err != nil {
		return nil, err
	}
	res.SimpleExpression = se

	for {
		if p.CurrentToken().Is(RelOpPredicator) {
			op := p.CurrentToken().Value()
			p.NextToken()
			se, err := p.ParseSimpleExpression()
			if err != nil {
				return nil, err
			}
			res.RelOpSimpleExpressions = append(
				res.RelOpSimpleExpressions,
				&ast.RelOpSimpleExpression{RelOp: op, SimpleExpression: se},
			)
		} else {
			break
		}
	}
	return res, nil
}

var RelOpPredicator = token.Some(
	token.Symbol('>'),
	token.Symbol('<'),
	token.SpecialSymbol.HasText("<="),
	token.SpecialSymbol.HasText(">="),
	token.Symbol('='),
	token.SpecialSymbol.HasText("<>"),
	token.ReservedWord.HasKeyword("IN"),
	token.ReservedWord.HasKeyword("IS"),
)

func (p *Parser) ParseSimpleExpression() (*ast.SimpleExpression, error) {
	res := &ast.SimpleExpression{}
	if p.CurrentToken().Is(UnaryOpPredicator) {
		s := p.CurrentToken().Value()
		res.UnaryOp = &s
		p.NextToken()
	}

	tm, err := p.ParseTerm()
	if err != nil {
		return nil, err
	}
	res.Term = tm

	for {
		if p.CurrentToken().Is(AddOpPredicator) {
			op := p.CurrentToken().Value()
			p.NextToken()
			tm, err := p.ParseTerm()
			if err != nil {
				return nil, err
			}
			res.AddOpTerms = append(
				res.AddOpTerms,
				&ast.AddOpTerm{AddOp: op, Term: tm},
			)
		} else {
			break
		}
	}
	return res, nil
}

var UnaryOpPredicator = token.Some(
	token.Symbol('+'),
	token.Symbol('-'),
)

var AddOpPredicator = token.Some(
	token.Symbol('+'),
	token.Symbol('-'),
	token.ReservedWord.HasKeyword("OR"),
	token.ReservedWord.HasKeyword("XOR"),
)

func (p *Parser) ParseTerm() (*ast.Term, error) {
	fac, err := p.ParseFactor()
	if err != nil {
		return nil, err
	}
	res := &ast.Term{Factor: fac}
	for {
		if p.CurrentToken().Is(MulOpPredicator) {
			op := p.CurrentToken().Value()
			p.NextToken()
			fac, err := p.ParseFactor()
			if err != nil {
				return nil, err
			}
			res.MulOpFactors = append(
				res.MulOpFactors,
				&ast.MulOpFactor{MulOp: op, Factor: fac},
			)
		} else {
			break
		}
	}
	return res, nil
}

var MulOpPredicator = token.Some(
	token.Symbol('*'),
	token.Symbol('/'),
	token.ReservedWord.HasKeyword("DIV"),
	token.ReservedWord.HasKeyword("MOD"),
	token.ReservedWord.HasKeyword("AND"),
	token.ReservedWord.HasKeyword("SHL"),
	token.ReservedWord.HasKeyword("SHR"),
)

func (p *Parser) ParseFactor() (ast.Factor, error) {
	t0 := p.CurrentToken()
	t0Value := t0.Value()
	if t0.Is(token.SpecialSymbol) {
		switch t0Value {
		case "@":
			p.NextToken()
			d, err := p.ParseDesignator()
			if err != nil {
				return nil, err
			}
			return &ast.Address{Designator: d}, nil
		case "[":
			set, err := p.ParseSetConstructor()
			if err != nil {
				return nil, err
			}
			p.NextToken()
			return set, nil
		case "(":
			p.NextToken()
			expr, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			if _, err := p.Current(token.Symbol(')')); err != nil {
				return nil, err
			}
			p.NextToken()
			return &ast.Parentheses{Expression: expr}, nil
		}
	} else if t0.Is(token.ReservedWord) {
		switch t0Value {
		case "NIL":
			p.NextToken()
			return &ast.Nil{}, nil
		case "NOT":
			p.NextToken()
			f, err := p.ParseFactor()
			if err != nil {
				return nil, err
			}
			return &ast.Not{Factor: f}, nil
		}
	} else if ast.IsManifestConstant(t0Value) {
		return p.ParseManifestConstant(t0, true)
	} else if t0.Is(token.CharacterString) {
		return p.ParseStringFactor(t0, true)
	} else if t0.Is(token.Some(token.NumeralInt, token.NumeralReal)) {
		return p.ParseNumberFactor(t0, true)
	} else if t0.Is(token.Some(token.Identifier)) {

		decl := p.context.Get(t0Value)
		if decl != nil {
			if _, ok := decl.Node.(*ast.TypeDecl); ok {
				rollback := p.RollbackPoint()
				typeId := ast.NewTypeId(ast.NewIdent(t0), decl)
				p.NextToken()
				if p.CurrentToken().Is(token.Symbol('(')) {
					p.NextToken()
					expr, err := p.ParseExpression()
					if err != nil {
						return nil, err
					}
					if _, err := p.Current(token.Symbol(')')); err != nil {
						return nil, err
					}
					p.NextToken()
					return &ast.TypeCast{TypeId: typeId, Expression: expr}, nil
				} else {
					rollback()
				}
			}
		}

		d, err := p.ParseDesignator()
		if err != nil {
			return nil, err
		}
		res := &ast.DesignatorFactor{
			Designator: d,
		}
		if p.CurrentToken().Is(token.Symbol('(')) {
			p.NextToken()
			exprList, err := p.ParseExprList(token.Symbol(')'))
			if err != nil {
				return nil, err
			}
			p.NextToken()
			res.ExprList = exprList
		}
		return res, nil
	}

	return nil, p.TokenErrorf("unexpected token %s", t0)
}

func (p *Parser) ParseManifestConstant(t *token.Token, skipTypeCheck bool) (*ast.ValueFactor, error) {
	if skipTypeCheck || ast.IsManifestConstant(t.Value()) {
		p.NextToken()
		return &ast.ValueFactor{Value: t.Value()}, nil
	} else {
		return nil, p.TokenErrorf("unexpected token %s for ValueFactor", t)
	}
}

func (p *Parser) ParseStringFactor(t *token.Token, skipTypeCheck bool) (*ast.StringFactor, error) {
	if skipTypeCheck || t.Is(token.CharacterString) {
		p.NextToken()
		return &ast.StringFactor{Value: t.Value()}, nil
	} else {
		return nil, p.TokenErrorf("unexpected token %s for StringFactor", t)
	}
}

func (p *Parser) ParseNumberFactor(t *token.Token, skipTypeCheck bool) (*ast.NumberFactor, error) {
	if skipTypeCheck || t.Is(token.CharacterString) {
		p.NextToken()
		return &ast.NumberFactor{Value: t.Value()}, nil
	} else {
		return nil, p.TokenErrorf("unexpected token %s for NumberFactor", t)
	}
}

func (p *Parser) ParseDesignator() (*ast.Designator, error) {
	if _, err := p.Current(token.Some(token.Identifier)); err != nil {
		return nil, err
	}
	res := &ast.Designator{}

	qualId, err := p.ParseQualId()
	if err != nil {
		return nil, err
	}
	res.QualId = qualId

	for {
		if _, err := p.Current(token.SpecialSymbol); err != nil {
			break
		}
		var item ast.DesignatorItem
		switch p.CurrentToken().Value() {
		case ".":
			t, err := p.Next(token.Identifier)
			if err != nil {
				return nil, err
			}
			item = ast.NewDesignatorItemIdent(t)
		case "[":
			p.NextToken()
			exprList, err := p.ParseExprList(token.Symbol(']'))
			if err != nil {
				return nil, err
			}
			item = ast.DesignatorItemExprList(exprList)
		case "^":
			item = &ast.DesignatorItemDereference{}
		}
		if item == nil {
			break
		}
		res.Items = append(res.Items, item)
		p.NextToken()
	}
	return res, nil
}

func (p *Parser) ParseSetConstructor() (*ast.SetConstructor, error) {
	if _, err := p.Current(token.Symbol('[')); err != nil {
		return nil, err
	}
	p.NextToken()
	res := &ast.SetConstructor{}
	if err := p.Until(token.Symbol(']'), token.Symbol(','), func() error {
		element, err := p.ParseSetElement()
		if err != nil {
			return err
		}
		res.SetElements = append(res.SetElements, element)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseSetElement() (*ast.SetElement, error) {
	res := &ast.SetElement{}
	expr1, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	res.Expression = expr1
	if p.CurrentToken().Is(token.SpecialSymbol.HasText("..")) {
		p.NextToken()
		expr2, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res.SubRangeEnd = expr2
	}
	return res, nil
}
