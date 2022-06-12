package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseStringType(required bool) (ast.StringType, error) {
	t := p.CurrentToken()
	decl := ast.EmbeddedTypeDecl(ast.EtkStringType, t.Value())
	if decl != nil {
		p.NextToken()
		return ast.NewTypeId(p.NewIdent(t), decl), nil
	} else if required {
		return nil, p.TokenErrorf("Unsupported token %s for StringType", t)
	} else {
		return nil, nil
	}
}

// This method parses just STRING not ANSISTRING nor WIDESTRING
func (p *Parser) ParseStringOfStringType() (ast.StringType, error) {
	t1 := p.CurrentToken()
	decl := ast.EmbeddedTypeDecl(ast.EtkStringType, t1.Value())
	if decl == nil {
		return nil, nil
	}
	t2 := p.NextToken()
	if t2.Is(token.Symbol('[')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(token.Symbol(']')); err != nil {
			return nil, err
		}
		p.NextToken()
		return ast.NewFixedStringType(p.NewIdent(t1), expr), nil
	} else {
		return ast.NewTypeId(p.NewIdent(t1), decl), nil
	}
}
