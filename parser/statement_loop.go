package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseRepeatStmt() (*ast.RepeatStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("REPEAT")); err != nil {
		return nil, err
	}
	p.NextToken()
	terminator := token.ReservedWord.HasKeyword("UNTIL")
	stmeList, err := p.ParseStmtList(terminator)
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(terminator); err != nil {
		return nil, err
	}
	p.NextToken()
	condition, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &ast.RepeatStmt{
		StmtList:  stmeList,
		Condition: condition,
	}, nil
}
