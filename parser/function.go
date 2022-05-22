package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseProcedureDeclSection() (*ast.FunctionDecl, error) {
	defer p.StackContext()()

	var functionHeading *ast.FunctionHeading
	switch p.CurrentToken().Value() {
	case "PROCEDURE":
		var err error
		functionHeading, err = p.ParseProcedureHeading()
		if err != nil {
			return nil, err
		}
	case "FUNCTION":
		var err error
		functionHeading, err = p.ParseFunctionHeading()
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	if _, err := p.Current(token.Symbol(';')); err != nil {
		return nil, err
	}
	res := &ast.FunctionDecl{FunctionHeading: functionHeading}
	if err := p.context.Set(res); err != nil {
		return nil, err
	}

	p.NextToken()
	if p.CurrentToken().Is(token.Directive) {
		directives, opts, err := p.ParseFunctionDirectives()
		if err != nil {
			return nil, err
		}
		res.Directives = directives
		res.ExternalOptions = opts
	}
	// TODO PortabilityDirective
	block, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}
	res.Block = block
	return res, nil
}
