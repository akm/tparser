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

func (p *Parser) ParseCaseStmt() (*ast.CaseStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("CASE")); err != nil {
		return nil, err
	}
	p.NextToken()
	expression, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("OF")); err != nil {
		return nil, err
	}
	p.NextToken()
	selectors, err := p.ParseCaseSelectors()
	if err != nil {
		return nil, err
	}
	terminator := token.ReservedWord.HasKeyword("END")
	var elseStmtList ast.StmtList
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("ELSE")) {
		p.NextToken()
		var err error
		elseStmtList, err = p.ParseStmtList(terminator)
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.Current(terminator); err != nil {
		return nil, err
	}
	p.NextToken()
	return &ast.CaseStmt{
		Expression: expression,
		Selectors:  selectors,
		Else:       elseStmtList,
	}, nil
}

func (p *Parser) ParseCaseSelectors() (ast.CaseSelectors, error) {
	r := ast.CaseSelectors{}
	for {
		selector, err := p.ParseCaseSelector()
		if err != nil {
			return nil, err
		}
		r = append(r, selector)
		if _, err := p.Current(token.Symbol(';')); err != nil {
			return nil, err
		}
		p.NextToken()
		if p.CurrentToken().Is(token.Some(
			token.ReservedWord.HasKeyword("ELSE"),
			token.ReservedWord.HasKeyword("END"),
		)) {
			break
		}
	}
	return r, nil
}

func (p *Parser) ParseCaseSelector() (*ast.CaseSelector, error) {
	labels, err := p.ParseCaseLabels()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.Symbol(':')); err != nil {
		return nil, err
	}
	p.NextToken()
	statement, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}
	return &ast.CaseSelector{
		Labels:    labels,
		Statement: statement,
	}, nil
}

func (p *Parser) ParseCaseLabels() (ast.CaseLabels, error) {
	labels := ast.CaseLabels{}
	for {
		label, err := p.ParseCaseLabel()
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
		if p.CurrentToken().Is(token.Symbol(',')) {
			p.NextToken()
		} else {
			break
		}
	}
	return labels, nil
}

func (p *Parser) ParseCaseLabel() (*ast.CaseLabel, error) {
	expr1, err := p.ParseConstExpr()
	if err != nil {
		return nil, err
	}
	var expr2 *ast.Expression
	if p.CurrentToken().Is(token.Symbol('.', '.')) {
		p.NextToken()
		var err error
		expr2, err = p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
	}
	return &ast.CaseLabel{
		ConstExpr:      expr1,
		ExtraConstExpr: expr2,
	}, nil
}
