package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseConstSection(required bool) (ast.ConstSection, error) {
	kw := token.ReservedWord.HasKeyword("CONST")
	if required {
		if _, err := p.Current(kw); err != nil {
			return nil, err
		}
	} else {
		if !p.CurrentToken().Is(kw) {
			return nil, nil
		}
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
	res.Ident = p.NewIdent(ident)
	if err := p.context.Set(res); err != nil {
		return nil, err
	}

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
