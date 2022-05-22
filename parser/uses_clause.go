package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseUsesClause() (ast.UsesClause, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("USES")); err != nil {
		return nil, err
	}
	r := ast.UsesClause{}
	p.NextToken()
	if err := p.Until(token.Symbol(';'), token.Symbol(','), func() error {
		t, err := p.Current(token.Identifier)
		if err != nil {
			return err
		}
		item := &ast.UsesClauseItem{Ident: p.NewIdent(t)}
		t2 := p.NextToken()
		if t2.Is(token.ReservedWord.HasKeyword("IN")) {
			t := p.NextToken()
			strFactor, err := p.ParseStringFactor(t, false)
			if err != nil {
				return err
			}
			item.Path = &strFactor.Value
		}
		r = append(r, item)
		return nil
	}); err != nil {
		return nil, err
	}
	for _, i := range r {
		if err := p.context.Set(i); err != nil {
			return nil, err
		}
	}
	return r, nil
}
