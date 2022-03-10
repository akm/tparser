package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseIdentClause() (*ast.IdentList, error) {
	res := ast.IdentList{}
	err := p.Until(token.Symbol(';'), token.Symbol(','), func() error {
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
