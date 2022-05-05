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
	ident, err := p.Current(token.Some(token.Identifier, token.Directive))
	if err != nil {
		return nil, err
	}
	res.Ident = ast.NewIdent(ident)
	p.context.DeclarationMap.SetDecl(res)

	p.NextToken()
	if p.CurrentToken().Is(token.Symbol(':')) {
		p.NextToken()
		typ, err := p.ParseType() // TODO TypeId, Array, Procedure or Pointer
		if err != nil {
			return nil, err
		}
		res.Type = typ
	}
	if _, err := p.Current(token.Symbol('=')); err != nil {
		return nil, err
	}

	p.NextToken()
	expr, err := p.ParseConstExpr()
	if err != nil {
		return nil, err
	}
	res.ConstExpr = expr

	return res, nil
}

func (p *Parser) ParseConstExpr() (*ast.ConstExpr, error) {
	// TODO Allow ConstExpr only
	return p.ParseExpression()
}
