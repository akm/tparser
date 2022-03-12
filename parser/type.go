package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) ParseTypeSection() (ast.TypeSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("TYPE")); err != nil {
		return nil, err
	}
	p.NextToken()
	res := ast.TypeSection{}
	for {
		decl, err := p.ParseTypeDecl()
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

func (p *Parser) ParseTypeDecl() (*ast.TypeDecl, error) {
	res := &ast.TypeDecl{}
	ident, err := p.Current(token.Identifier)
	if err != nil {
		return nil, err
	}
	res.Ident = ast.Ident(ident.Value())
	if _, err := p.Next(token.Symbol('=')); err != nil {
		return nil, err
	}
	t := p.NextToken()
	if t.Is(token.ReservedWord.HasKeyword("TYPE")) {
		// ignore
		t = p.NextToken()
	}
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ

	{
		t := p.CurrentToken()
		if t.Is(token.Symbol(';')) {
			return res, nil
		}
		if t := p.NextToken(); t.Is(token.PortabilityDirective) {
			portabilityDirective := ast.PortabilityDirective(t.Value())
			res.PortabilityDirective = &portabilityDirective
			p.NextToken()
		}
	}
	return res, nil
}

func (p *Parser) ParseType() (ast.Type, error) {
	t1 := p.CurrentToken()
	switch t1.Type {
	case token.Identifier:
		return p.ParseTypeId()
	}
	return nil, errors.Errorf("Unsupported Type token %+v", t1)
}

func (p *Parser) ParseTypeId() (ast.Type, error) {
	t1 := p.CurrentToken()
	part1 := t1.Value()
	t2 := p.NextToken()
	if t2.Is(token.Symbol('.')) {
		t3, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		unitId := ast.UnitId(part1)
		return &ast.TypeId{
			UnitId: &unitId,
			Ident:  ast.Ident(t3.Value()),
		}, nil
	} else {
		return &ast.TypeId{
			Ident: ast.Ident(part1),
		}, nil
	}
}
