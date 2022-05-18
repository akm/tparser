package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) NewIdent(t *token.Token) *ast.Ident {
	return ast.NewIdent(t)
}

func (p *Parser) NewIdentRef(t *token.Token) *ast.IdentRef {
	return ast.NewIdentRef(p.NewIdent(t), p.context.Get(t.RawString()))
}
