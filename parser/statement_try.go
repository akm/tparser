package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/log"
	"github.com/akm/tparser/token"
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
	} else if p.CurrentToken().Is(token.ReservedWord.HasKeyword("FINALLY")) {
		p.NextToken()
		finallyStmtList, err := p.ParseStmtList(token.ReservedWord.HasKeyword("END"))
		if err != nil {
			return nil, err
		}
		p.NextToken()
		return &ast.TryFinallyStmt{
			Statements1: stmtList,
			Statements2: finallyStmtList,
		}, nil
	} else {
		return nil, p.TokenErrorf("expected 'except' or 'finally' but got %s", p.CurrentToken())
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

	res := &ast.ExceptionBlock{Handlers: handlers}

	if p.CurrentToken().Is(kwEnd) {
		p.NextToken()
		return res, nil
	}

	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("ELSE")) {
		p.NextToken()
		statements, err := p.ParseStmtList(kwEnd)
		if err != nil {
			return nil, err
		}
		res.Else = statements
	}

	if p.CurrentToken().Is(token.Symbol(';')) {
		p.NextToken()
	}

	if _, err := p.Current(kwEnd); err != nil {
		return nil, err
	}
	p.NextToken()
	return res, nil
}

func (p *Parser) ParseExceptionBlockHandlers() (ast.ExceptionBlockHandlers, error) {
	res := ast.ExceptionBlockHandlers{}
	for {
		handler, err := p.ParseExceptionBlockHandler()
		if err != nil {
			return nil, err
		}
		res = append(res, handler)

		if p.CurrentToken().Is(token.Symbol(';')) {
			p.NextToken()
		}

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

	// p.logger.Printf("p.context.DeclMap.Keys(): %+v\n", p.context.DeclMap.Keys())

	if decl := p.context.Get(t.RawString()); decl != nil {
		log.Printf("decl: %+v\n", *decl)
		log.Printf("decl.Node: %+v\n", decl.Node)
		if _, ok := decl.Node.(*ast.TypeDecl); ok {
			log.Printf("decl is Node\n")
			hasIdent = false
		} else {
			log.Printf("decl is NOT Node\n")
		}
	}
	decl := &ast.ExceptionBlockHandlerDecl{}
	if hasIdent {
		decl.Ident = p.NewIdent(t)
		if _, err := p.Next(token.Symbol(':')); err != nil {
			return nil, err
		}
		p.NextToken()
	}
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	decl.Type = typ
	if decl.Ident != nil {
		if err := p.context.Set(decl); err != nil {
			return nil, err
		}
	}

	res := &ast.ExceptionBlockHandler{Decl: decl}

	if _, err := p.Current(token.UpperCase("DO")); err != nil {
		return nil, err
	}
	p.NextToken()

	statement, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}
	res.Statement = statement
	return res, nil
}
