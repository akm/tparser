package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) NewIdent(t *token.Token) *ast.Ident {
	return ast.NewIdent(t)
}
