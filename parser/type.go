package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) ParseTypeSection() (ast.TypeSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("TYPE")); err != nil {
		return nil, err
	}
	p.NextToken()
	res := ast.TypeSection{}
	for {
		decl, err := p.ParseTypeDecl()
		if err != nil {
			return nil, err
		}
		if _, err := p.Current(token.Symbol(';')); err != nil {
			return nil, err
		}
		res = append(res, decl)
		t := p.NextToken()
		if t.Is(token.ReservedWord) || t.Is(token.EOF) {
			break
		}
	}
	return res, nil
}

func (p *Parser) ParseTypeDecl() (*ast.TypeDecl, error) {
	res := &ast.TypeDecl{}
	ident, err := p.Current(token.Identifier)
	if err != nil {
		return nil, err
	}
	res.Ident = ast.Ident(ident.Value())
	if _, err := p.Next(token.Symbol('=')); err != nil {
		return nil, err
	}
	t := p.NextToken()
	if t.Is(token.ReservedWord.HasKeyword("TYPE")) {
		// ignore
		t = p.NextToken()
	}
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ

	{
		t := p.CurrentToken()
		if t.Is(token.Symbol(';')) {
			return res, nil
		}
		if t := p.NextToken(); t.Is(token.PortabilityDirective) {
			portabilityDirective := ast.PortabilityDirective(t.Value())
			res.PortabilityDirective = &portabilityDirective
			p.NextToken()
		}
	}
	return res, nil
}

func (p *Parser) ParseType() (ast.Type, error) {
	t1 := p.CurrentToken()
	switch t1.Type {
	case token.SpecialSymbol:
		if t1.Is(token.Symbol('(')) {
			return p.ParseEnumeratedType()
		}
	case token.Identifier:
		if typ, err := p.ParseNamedType(); err != nil {
			return nil, err
		} else if typ != nil {
			return typ, nil
		}
		return p.ParseTypeIdOrSubrangeType()
	case token.NumeralInt, token.NumeralReal, token.CharacterString:
		return p.ParseConstSubrageType()
	case token.ReservedWord:
		return p.ParseStringOfStringType()
	}
	return nil, errors.Errorf("Unsupported Type token %+v", t1)
}

func (p *Parser) ParseNamedType() (ast.Type, error) {
	t1 := p.CurrentToken()
	name := t1.Value()
	if ast.IsRealTypeName(name) {
		return &ast.RealType{Name: ast.Ident(name)}, nil
	} else if ast.IsOrdIdentName(name) {
		return &ast.OrdIdent{Name: ast.Ident(name)}, nil
	} else if ast.IsStringTypeName(name) {
		return &ast.StringType{Name: name}, nil
	} else {
		return nil, nil
	}
}

func (p *Parser) ParseTypeIdOrSubrangeType() (ast.Type, error) {
	t1 := p.CurrentToken()
	part1 := t1.Value()
	t2 := p.NextToken()
	if t2.Is(token.Value("..")) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		return &ast.SubrangeType{
			Low:  *ast.NewConstExpr(part1),
			High: *expr,
		}, nil
	} else if t2.Is(token.Symbol('.')) {
		t3, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		unitId := ast.UnitId(part1)
		return &ast.TypeId{
			UnitId: &unitId,
			Ident:  ast.Ident(t3.Value()),
		}, nil
	} else {
		return &ast.TypeId{
			Ident: ast.Ident(part1),
		}, nil
	}
}

func (p *Parser) ParseEnumeratedType() (ast.EnumeratedType, error) {
	res := ast.EnumeratedType{}
	for {
		element, err := p.ParseEnumeratedTypeElement()
		if err != nil {
			return nil, err
		}
		res = append(res, element)
		t := p.NextToken()
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
	// TODO parse ConstExpr if exists
	return &ast.EnumeratedTypeElement{Ident: ast.Ident(ident.Value())}, nil
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
		// TODO parse ConstExpr
		t3 := p.NextToken()
		if _, err := p.Next(token.Symbol(']')); err != nil {
			return nil, err
		}

		l := t3.Value()
		return &ast.StringType{Name: "STRING", Length: &l}, nil
	} else {
		return nil, errors.Errorf("unexpected token %s", t2)
	}
}
