package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseConstSection() (ast.ConstSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("CONST")); err != nil {
		return nil, err
	}
	p.NextToken()
	res := ast.ConstSection{}
	for {
		decl, err := p.ParseConstantDecl()
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(token.Symbol(';')); err != nil {
			return nil, err
		}
		res = append(res, decl)
		t := p.NextToken()
		if t.Is(token.ReservedWord) || t.Is(token.EOF) {
			break
		}
	}
	return res, nil
}

func (p *Parser) ParseConstantDecl() (*ast.ConstantDecl, error) {
	res := &ast.ConstantDecl{}
	ident, err := p.Current(token.Identifier)
	if err != nil {
		return nil, err
	}
	res.Ident = ast.Ident(ident.Value())

	p.NextToken()
	if p.CurrentToken().Is(token.Symbol(':')) {
		p.NextToken()
		typ, err := p.ParseType() // TODO TypeId, Array, Procedure or Pointer
		if err != nil {
			return nil, err
		}
		res.Type = typ
		p.NextToken()
	}
	if _, err := p.Current(token.Symbol('=')); err != nil {
		return nil, err
	}
	t := p.NextToken()
	res.ConstExpr.Value = t.Value()
	p.NextToken()

	return res, nil
}