package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

var tryStmtTerminator = token.Some(
	token.ReservedWord.HasKeyword("EXCEPT"),
	token.ReservedWord.HasKeyword("FINALLY"),
)

func (p *Parser) ParseTryStmt() (ast.TryStmt, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("TRY")); err != nil {
		return nil, err
	}
	p.NextToken()
	stmtList, err := p.ParseStmtList(tryStmtTerminator)
	if err != nil {
		return nil, err
	}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("EXCEPT")) {
		p.NextToken()
		exceptionBlock, err := p.ParseExceptionBlock()
		if err != nil {
			return nil, err
		}
		return &ast.TryExceptStmt{
			Statements:     stmtList,
			ExceptionBlock: exceptionBlock,
		}, nil
	} else {
		return nil, errors.Errorf("expected 'except' but got %s", p.CurrentToken().RawString())
	}
}

func (p *Parser) ParseExceptionBlock() (*ast.ExceptionBlock, error) {
	kwEnd := token.ReservedWord.HasKeyword("END")
	// ON is NOT a reserved word
	// If there is no "ON" at the head, then the exception block is else statements only
	if !p.CurrentToken().Is(token.UpperCase("ON")) {
		statements, err := p.ParseStmtList(kwEnd)
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(kwEnd); err != nil {
			return nil, err
		}
		p.NextToken()
		return &ast.ExceptionBlock{Else: statements}, nil
	}
	handlers, err := p.ParseExceptionBlockHandlers()
	if err != nil {
		return nil, err
	}
	if p.CurrentToken().Is(kwEnd) {
		p.NextToken()
		return &ast.ExceptionBlock{Handlers: handlers}, nil
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("ELSE")); err != nil {
		return nil, err
	}

	statements, err := p.ParseStmtList(kwEnd)
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(kwEnd); err != nil {
		return nil, err
	}
	p.NextToken()
	return &ast.ExceptionBlock{Handlers: handlers, Else: statements}, nil
}

func (p *Parser) ParseExceptionBlockHandlers() (ast.ExceptionBlockHandlers, error) {
	res := ast.ExceptionBlockHandlers{}
	for {
		handler, err := p.ParseExceptionBlockHandler()
		if err != nil {
			return nil, err
		}
		res = append(res, handler)
		// ON is NOT a reserved word
		if !p.CurrentToken().Is(token.UpperCase("ON")) {
			break
		}
	}
	return res, nil
}

func (p *Parser) ParseExceptionBlockHandler() (*ast.ExceptionBlockHandler, error) {
	// ON is NOT a reserved word
	if _, err := p.Current(token.UpperCase("ON")); err != nil {
		return nil, err
	}
	t := p.NextToken()

	hasIdent := true
	if decl := p.context.DeclarationMap.Get(t.RawString()); decl != nil {
		if _, ok := decl.Node.(ast.Type); ok {
			hasIdent = false
		}
	}
	res := &ast.ExceptionBlockHandler{}
	if hasIdent {
		res.Ident = ast.NewIdent(t)
		if _, err := p.Next(token.Symbol(':')); err != nil {
			return nil, err
		}
		p.NextToken()
	}
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ
	statement, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}
	res.Statement = statement
	return res, nil
}
