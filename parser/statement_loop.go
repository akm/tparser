package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
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

func (p *Parser) ParseWhileStmt() (*ast.WhileStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("WHILE")); err != nil {
		return nil, err
	}
	p.NextToken()
	condition, err := p.ParseExpression()
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
	return &ast.WhileStmt{
		Condition: condition,
		Statement: statement,
	}, nil
}

func (p *Parser) ParseForStmt() (*ast.ForStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("FOR")); err != nil {
		return nil, err
	}
	p.NextToken()
	qualId, err := p.ParseQualId()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.Symbol(':', '=')); err != nil {
		return nil, err
	}
	p.NextToken()
	initial, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	var down bool
	switch p.CurrentToken().Value() {
	case "TO":
		down = false
	case "DOWNTO":
		down = true
	default:
		return nil, errors.Errorf("expected TO or DOWNTO, but got %s", p.CurrentToken().Value())
	}
	p.NextToken()
	terminal, err := p.ParseExpression()
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
	return &ast.ForStmt{
		QualId:    qualId,
		Initial:   initial,
		Terminal:  terminal,
		Down:      down,
		Statement: statement,
	}, nil
}
