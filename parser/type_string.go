package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

var stringType = token.PredicatorBy("StringType", ast.IsStringTypeName)

func (p *Parser) ParseStringType(required bool) (*ast.StringType, error) {
	t := p.CurrentToken()
	if t.Is(stringType) {
		p.NextToken()
		return &ast.StringType{Name: t.Value()}, nil
	} else if required {
		return nil, errors.Errorf("Unsupported token %+v for StringType", t)
	} else {
		return nil, nil
	}
}

// This method parses just STRING not ANSISTRING nor WIDESTRING
func (p *Parser) ParseStringOfStringType() (*ast.StringType, error) {
	t1 := p.CurrentToken()
	if !t1.Is(token.Value("STRING")) {
		return nil, nil
	}
	t2 := p.NextToken()
	if t2.Is(token.Symbol(';')) || t2.Is(token.EOF) {
		return &ast.StringType{Name: t1.Value()}, nil
	} else if t2.Is(token.Symbol('[')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(token.Symbol(']')); err != nil {
			return nil, err
		}
		p.NextToken()
		return &ast.StringType{Name: "STRING", Length: expr}, nil
	} else {
		return nil, errors.Errorf("unexpected token %s", t2)
	}
}
