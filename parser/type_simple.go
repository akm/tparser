package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

var realType = token.PredicatorBy("RealType", ast.IsRealTypeName)

func (p *Parser) ParseRealType(required bool) (*ast.RealType, error) {
	t := p.CurrentToken()
	if t.Is(realType) {
		p.NextToken()
		return &ast.RealType{Name: *ast.NewIdent(t)}, nil
	} else if required {
		return nil, errors.Errorf("Unsupported token %+v for RealType", t)
	} else {
		return nil, nil
	}
}

var ordIdent = token.PredicatorBy("OrdIdent", ast.IsOrdIdentName)

func (p *Parser) ParseOrdIdent(required bool) (*ast.OrdIdent, error) {
	t := p.CurrentToken()
	if t.Is(ordIdent) {
		p.NextToken()
		return &ast.OrdIdent{Name: *ast.NewIdent(t)}, nil
	} else if required {
		return nil, errors.Errorf("Unsupported token %+v for OrdIdent", t)
	} else {
		return nil, nil
	}
}

func (p *Parser) parseSubrangeTypeForIdentifier(required bool) (*ast.SubrangeType, error) {
	rollback := p.RollbackPoint()
	t1 := p.CurrentToken()
	t2 := p.NextToken()

	if t2.Is(token.Value("..")) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		return &ast.SubrangeType{
			Low:  *ast.NewConstExpr(t1),
			High: *expr,
		}, nil
	} else {
		defer rollback()
		if required {
			return nil, errors.Errorf("Unsupported token %+v, %+v for SubrangeType", t1, t2)
		} else {
			return nil, nil
		}
	}
}

func (p *Parser) ParseConstSubrageType() (*ast.SubrangeType, error) {
	lowExpr, err := p.ParseConstExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.Value("..")); err != nil {
		return nil, err
	}
	p.NextToken()
	highExpr, err := p.ParseConstExpr()
	if err != nil {
		return nil, err
	}
	return &ast.SubrangeType{
		Low:  *lowExpr,
		High: *highExpr,
	}, nil
}

func (p *Parser) ParseEnumeratedType() (ast.EnumeratedType, error) {
	res := ast.EnumeratedType{}
	for {
		element, err := p.ParseEnumeratedTypeElement()
		if err != nil {
			return nil, err
		}
		res = append(res, element)
		t := p.CurrentToken()
		if t.Is(token.Symbol(')')) {
			break
		} else if t.Is(token.Symbol(',')) {
			continue
		} else {
			return nil, errors.Errorf("Unsupported token %+v for EnumeratedType", t)
		}
	}
	return res, nil
}

func (p *Parser) ParseEnumeratedTypeElement() (*ast.EnumeratedTypeElement, error) {
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.EnumeratedTypeElement{Ident: *ast.NewIdent(ident)}
	p.NextToken()
	if p.CurrentToken().Is(token.Symbol('=')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		res.ConstExpr = expr
	}

	return res, nil
}
