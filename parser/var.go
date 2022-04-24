package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseVarSection() (ast.VarSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("VAR")); err != nil {
		return nil, err
	}
	p.NextToken()
	res := ast.VarSection{}
	for {
		decl, err := p.ParseVarDecl()
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

func (p *Parser) ParseVarDecl() (*ast.VarDecl, error) {
	res := &ast.VarDecl{}
	identList, err := p.ParseIdentList(':')
	if err != nil {
		return nil, err
	}
	res.IdentList = *identList

	p.NextToken()

	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ

	if p.CurrentToken().Is(token.UpperCase("ABSOLUTE")) {
		// TODO Support ConstExpr for absolute but no example found for this
		t, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		res.Absolute = ast.NewVarDeclAbsoluteIdent(t.Value())
		p.NextToken()
	}

	if p.CurrentToken().Is(token.Symbol('=')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		res.ConstExpr = expr
	}
	return res, nil
}

func (p *Parser) ParseThreadVarSection() (ast.ThreadVarSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("THREADVAR")); err != nil {
		return nil, err
	}
	p.NextToken()
	res := ast.ThreadVarSection{}
	for {
		decl, err := p.ParseThreadVarDecl()
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

func (p *Parser) ParseThreadVarDecl() (*ast.ThreadVarDecl, error) {
	res := &ast.ThreadVarDecl{}
	identList, err := p.ParseIdentList(':')
	if err != nil {
		return nil, err
	}
	res.IdentList = *identList

	p.NextToken()

	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ

	// p.NextToken()
	return res, nil
}
