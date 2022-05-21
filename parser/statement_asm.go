package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/log"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseAssemblerStatement() (*ast.AssemblerStatement, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("ASM")); err != nil {
		return nil, err
	}
	p.NextToken()

	log.Printf("ParseAssemblerStatement start")

	for {
		if p.NextToken().Is(token.ReservedWord.HasKeyword("END")) {
			p.NextToken()
			break
		}
	}

	log.Printf("ParseAssemblerStatement done")

	return &ast.AssemblerStatement{}, nil
}
