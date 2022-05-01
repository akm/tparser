package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseIdentList(terminator rune) (*ast.IdentList, error) {
	return p.ParseIdentListBy(token.Symbol(terminator))
}

func (p *Parser) ParseIdentListBy(terminatorPredicate token.Predicator) (*ast.IdentList, error) {
	res := ast.IdentList{}
	err := p.Until(terminatorPredicate, token.Symbol(','), func() error {
		t, err := p.Current(token.Identifier)
		if err != nil {
			return err
		}
		res = append(res, ast.NewIdent(t))
		p.NextToken()
		return nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(terminatorPredicate); err != nil {
		return nil, err
	}
	return &res, nil
}
