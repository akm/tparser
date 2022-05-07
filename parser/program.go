package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseProgram() (*ast.Program, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("PROGRAM")); err != nil {
		return nil, err
	}
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Program{Ident: ast.NewIdent(ident)}
	if _, err := p.Next(token.Symbol(';')); err != nil {
		return nil, err
	}
	p.NextToken()
	block, err := p.ParseProgramBlock()
	if err != nil {
		return nil, err
	}
	res.ProgramBlock = block
	if _, err := p.Current(token.Symbol('.')); err != nil {
		return nil, err
	}
	p.context.DeclarationMap.SetDecl(res)
	return res, nil
}

func (p *Parser) ParseProgramBlock() (*ast.ProgramBlock, error) {
	res := &ast.ProgramBlock{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		uses, err := p.ParseUsesClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = uses
		p.NextToken()
	}
	block, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}
	res.Block = block
	return res, nil
}
