package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseCompoundStmt(required bool) (*ast.CompoundStmt, error) {
	kw := token.ReservedWord.HasKeyword("BEGIN")
	if required {
		if _, err := p.Current(kw); err != nil {
			return nil, err
		}
	} else {
		if !p.CurrentToken().Is(kw) {
			return nil, nil
		}
	}
	p.NextToken()

	terminator := token.ReservedWord.HasKeyword("END")
	stmtList, err := p.ParseStmtList(terminator)
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(terminator); err != nil {
		return nil, err
	}
	p.NextToken()
	return &ast.CompoundStmt{StmtList: stmtList}, nil
}

func (p *Parser) ParseStmtList(terminator token.Predicator) (ast.StmtList, error) {
	res := ast.StmtList{}
	for {
		statement, err := p.ParseStatement()
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(token.Symbol(';')); err != nil {
			return nil, err
		}
		res = append(res, statement)
		p.NextToken()
		if p.CurrentToken().Is(terminator) {
			break
		}
	}
	return res, nil
}

func (p *Parser) ParseStatement() (*ast.Statement, error) {
	res := &ast.Statement{}
	rollback := p.RollbackPoint()
	labelId := p.CurrentToken()
	if p.NextToken().Is(token.Symbol(':')) {
		res.LabelId = ast.NewLabelId(ast.NewIdent(labelId))
		p.NextToken()
	} else {
		rollback()
	}

	t := p.CurrentToken()
	if t.Is(token.ReservedWord) {
		switch t.Value() {
		case "INHERITED":
			if stmt, err := p.ParseInheritedStmt(); err != nil {
				return nil, err
			} else {
				res.Body = stmt
				return res, nil
			}
		case "GOTO":
			if stmt, err := p.ParseGotoStatement(); err != nil {
				return nil, err
			} else {
				res.Body = stmt
				return res, nil
			}
		}
	}

	// TODO InheritedStatement
	// TODO GotoStatement
	// TODO CompoundStmt
	// TODO ConditionalStmt
	// TODO LoopStmt
	// TODO WithStmt
	// TODO TryExceptStmt
	// TODO TryFinallyStmt
	// TODO RaiseStmt
	// TODO AssemblerStmt

	if stmt, err := p.ParseDesignatorStatement(); err != nil {
		return nil, err
	} else {
		res.Body = stmt
	}
	return res, nil
}

func (p *Parser) ParseDesignatorStatement() (ast.DesignatorStatement, error) {
	designator, err := p.ParseDesignator()
	if err != nil {
		return nil, err
	}
	t := p.CurrentToken()
	if t.Value() == ":=" {
		p.NextToken()
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		return &ast.AssignStatement{Designator: designator, Expression: expr}, nil
	} else {
		res := &ast.CallStatement{Designator: designator}
		if t.Is(token.Symbol('(')) {
			p.NextToken()
			terminator := token.Symbol(')')
			exprList, err := p.ParseExprList(terminator)
			if err != nil {
				return nil, err
			}
			res.ExprList = exprList
			if _, err := p.Current(terminator); err != nil {
				return nil, err
			}
			p.NextToken()
		}
		return res, nil
	}
}

func (p *Parser) ParseInheritedStmt() (*ast.InheritedStatement, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("INHERITED")); err != nil {
		return nil, err
	}
	p.NextToken()
	return &ast.InheritedStatement{
		// Ref: (find callee ancestor method)
	}, nil
}

func (p *Parser) ParseGotoStatement() (*ast.GotoStatement, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("GOTO")); err != nil {
		return nil, err
	}
	t, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	d := p.context.DeclarationMap.Get(t.Value())
	p.NextToken()
	return &ast.GotoStatement{
		LabelId: ast.NewLabelId(ast.NewIdent(t)),
		Ref:     d,
	}, nil
}
