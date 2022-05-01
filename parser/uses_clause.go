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
		ref := &ast.UnitRef{Ident: ast.NewIdent(t)}
		t2 := p.NextToken()
		if t2.Is(token.ReservedWord.HasKeyword("IN")) {
			t := p.NextToken()
			strFactor, err := p.ParseStringFactor(t, false)
			if err != nil {
				return err
			}
			ref.Path = &strFactor.Value
		}
		r = append(r, ref)
		return nil
	}); err != nil {
		return nil, err
	}
	return r, nil
}
