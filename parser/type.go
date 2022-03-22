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
		p.NextToken()
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
		return p.ParseTypeForIdentifier()
	case token.NumeralInt, token.NumeralReal, token.CharacterString:
		return p.ParseConstSubrageType()
	case token.ReservedWord:
		return p.ParseStringOfStringType()
	}
	return nil, errors.Errorf("Unsupported Type token %+v", t1)
}

func (p *Parser) ParseTypeForIdentifier() (ast.Type, error) {
	if res, err := p.ParseRealType(false); err != nil {
		return nil, err
	} else if res != nil {
		return res, nil
	} else if res, err := p.ParseOrdIdent(false); err != nil {
		return nil, err
	} else if res != nil {
		return res, nil
	} else if res, err := p.ParseStringType(false); err != nil {
		return nil, err
	} else if res != nil {
		return res, nil
	} else {
		return p.ParseTypeIdOrSubrangeType()
	}
}

var realType = token.PredicatorBy("RealType", ast.IsRealTypeName)

func (p *Parser) ParseRealType(required bool) (*ast.RealType, error) {
	t := p.CurrentToken()
	if t.Is(realType) {
		p.NextToken()
		return &ast.RealType{Name: ast.Ident(t.Value())}, nil
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
		return &ast.OrdIdent{Name: ast.Ident(t.Value())}, nil
	} else if required {
		return nil, errors.Errorf("Unsupported token %+v for OrdIdent", t)
	} else {
		return nil, nil
	}
}

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

func (p *Parser) ParseTypeIdOrSubrangeType() (ast.Type, error) {
	t1 := p.CurrentToken()
	t2 := p.NextToken()
	if res, err := p.parseSubrangeTypeForIdentifier(t1, t2, false); err != nil {
		return nil, err
	} else if res != nil {
		return res, nil
	} else if res, err := p.parseTypeId(t1, t2); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// t1 must be identifier token
// t2 can be ".." or others
func (p *Parser) parseSubrangeTypeForIdentifier(t1, t2 *token.Token, required bool) (*ast.SubrangeType, error) {
	if t2.Is(token.Value("..")) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		return &ast.SubrangeType{
			Low:  *ast.NewConstExpr(t1.Value()),
			High: *expr,
		}, nil
	} else if required {
		return nil, errors.Errorf("Unsupported token %+v, %+v for SubrangeType", t1, t2)
	} else {
		return nil, nil
	}
}

// t1 must be identifier token
// t2 can be "." or others
func (p *Parser) parseTypeId(t1, t2 *token.Token) (*ast.TypeId, error) {
	if t2.Is(token.Symbol('.')) {
		t3, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		unitId := ast.UnitId(t1.Value())
		p.NextToken()
		return &ast.TypeId{
			UnitId: &unitId,
			Ident:  ast.Ident(t3.Value()),
		}, nil
	} else {
		return &ast.TypeId{
			Ident: ast.Ident(t1.Value()),
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
	res := &ast.EnumeratedTypeElement{Ident: ast.Ident(ident.Value())}
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
