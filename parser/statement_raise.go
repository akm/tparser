package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseRaiseStmt() (*ast.RaiseStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("RAISE")); err != nil {
		return nil, err
	}
	p.NextToken()

	res := &ast.RaiseStmt{}

	if p.CurrentToken().Is(token.Symbol(';')) {
		return res, nil
	}

	if !p.CurrentToken().Is(token.UpperCase("AT")) {
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res.Object = expr
	}

	if p.CurrentToken().Is(token.UpperCase("AT")) {
		p.NextToken()
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		res.Address = expr
	}

	return res, nil
}
