package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseWithStmt() (*ast.WithStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("WITH")); err != nil {
		return nil, err
	}
	p.NextToken()
	qualIds, err := p.ParseQualIds()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("DO")); err != nil {
		return nil, err
	}
	p.NextToken()
	statement, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}
	return &ast.WithStmt{
		Objects:   qualIds,
		Statement: statement,
	}, nil
}
