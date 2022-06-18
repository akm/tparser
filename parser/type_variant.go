package parser

import (
	"github.com/akm/tparser/ast"
)

func (p *Parser) ParseVariantType(required bool) (ast.VariantType, error) {
	t := p.CurrentToken()
	decl := ast.EmbeddedTypeDecl(ast.EtkVariantType, t.Value())
	if decl != nil {
		p.NextToken()
		return ast.NewTypeId(p.NewIdent(t), decl), nil
	} else if required {
		return nil, p.TokenErrorf("Unsupported token %s for StringType", t)
	} else {
		return nil, nil
	}
}
