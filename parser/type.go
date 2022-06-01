package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) ParseTypeSection(required bool) (ast.TypeSection, error) {
	kw := token.ReservedWord.HasKeyword("TYPE")
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
	res.Ident = p.NewIdent(ident)
	if _, err := p.Next(token.Symbol('=')); err != nil {
		return nil, err
	}
	t := p.NextToken()
	if t.Is(token.ReservedWord.HasKeyword("TYPE")) {
		// TODO res.Typed = true
		p.NextToken()
	}
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	res.Type = typ
	if err := p.context.Set(res); err != nil {
		return nil, err
	}

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
		switch t1.Value() {
		case "PACKED", "ARRAY":
			return p.ParseStrucType()
		default:
			return p.ParseStringOfStringType()
		}
	}
	return nil, p.TokenErrorf("Unsupported Type token %s", t1)
}

func (p *Parser) ParseTypeForIdentifier() (ast.Type, error) {
	if res, err := p.parseTypeIdWithUnit(); res != nil || err != nil {
		return res, err
	} else if res, err := p.ParseRealType(false); res != nil || err != nil {
		return res, err
	} else if res, err := p.ParseOrdIdent(false); res != nil || err != nil {
		return res, err
	} else if res, err := p.ParseStringType(false); res != nil || err != nil {
		return res, err
	} else if res, err := p.parseSubrangeTypeForIdentifier(false); res != nil || err != nil {
		return res, err
	} else {
		return p.parseTypeIdWithoutUnit()
	}
}

func (p *Parser) parseTypeIdWithUnit() (*ast.TypeId, error) {
	if !p.IsUnitIdentifier(p.CurrentToken()) {
		return nil, nil
	}
	unitId := ast.NewUnitId(p.CurrentToken())
	if _, err := p.Next(token.Symbol('.')); err != nil {
		return nil, err
	}
	t, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	return &ast.TypeId{UnitId: unitId, Ident: p.NewIdent(t)}, nil
}

func (p *Parser) parseTypeIdWithoutUnit() (*ast.TypeId, error) {
	ident := p.NewIdent(p.CurrentToken())
	p.NextToken()
	r := &ast.TypeId{Ident: ident}

	decl := p.context.Get(ident.Name)
	if decl == nil {
		p.Logf("%s is not declared", ident.Name)
	} else {
		r.Ref = decl
	}

	return r, nil
}
