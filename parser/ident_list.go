package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseIdentList(terminator rune) (*ast.IdentList, error) {
	res := ast.IdentList{}
	err := p.Until(token.Symbol(terminator), token.Symbol(','), func() error {
		t, err := p.Next(token.Identifier)
		if err != nil {
			return err
		}
		res = append(res, t.Value())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
