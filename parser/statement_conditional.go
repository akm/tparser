package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseIfStmt() (*ast.IfStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("IF")); err != nil {
		return nil, err
	}
	p.NextToken()
	condition, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("THEN")); err != nil {
		return nil, err
	}
	p.NextToken()
	thenStmt, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}
	var elseStmt *ast.Statement
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("ELSE")) {
		p.NextToken()
		var err error
		elseStmt, err = p.ParseStatement()
		if err != nil {
			return nil, err
		}
	}
	return &ast.IfStmt{
		Condition: condition,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil
}
