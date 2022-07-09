package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseProcedureType() (*ast.ProcedureType, error) {
	res := &ast.ProcedureType{}

	t0 := p.CurrentToken()
	switch t0.Value() {
	case "FUNCTION":
		res.FunctionType = ast.FtFunction
	case "PROCEDURE":
		res.FunctionType = ast.FtProcedure
	default:
		return nil, p.TokenErrorf("expects FUNCTION or PROCEDURE, but got %s (%s)", t0, string(t0.Raw()))
	}

	defer p.context.StackDeclMap()()

	t := p.NextToken()
	if t.Is(token.Symbol('(')) {
		formalParameters, err := p.ParseFormalParameters('(', ')')
		if err != nil {
			return nil, err
		}
		res.FormalParameters = formalParameters
	}
	if res.FunctionType == ast.FtFunction {
		if _, err := p.Current(token.Symbol(':')); err != nil {
			return nil, err
		}
		p.NextToken()
		typ, err := p.ParseTypeId()
		if err != nil {
			return nil, err
		}
		res.ReturnType = typ
	}

	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("OF")) {
		p.NextToken()
		if _, err := p.Current(token.ReservedWord.HasKeyword("OBJECT")); err != nil {
			return nil, err
		}
		res.OfObject = true
		p.NextToken()
	}

	return res, nil
}
