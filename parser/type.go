package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
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
		{
			t := p.CurrentToken()
			if t.Is(token.ReservedWord) || t.Is(token.EOF) {
				res = append(res, decl)
				break
			}
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
	if old := p.context.Get(res.Ident.Name); old != nil {
		typeDecl, ok := old.Node.(*ast.TypeDecl)
		if !ok {
			return nil, errors.Errorf("expects type declaration but was %T", old.Node)
		}
		fwd, ok := typeDecl.Type.(ast.ForwardDeclaration)
		if !ok {
			return nil, errors.Errorf("expects forward declaration but was %T", typeDecl.Type)
		}
		if err := fwd.SetActualType(typ); err != nil {
			return nil, err
		}
		for _, decl := range res.ToDeclarations() {
			p.context.Overwrite(res.Ident.Name, decl)
		}
	} else {
		if err := p.context.Set(res); err != nil {
			return nil, err
		}
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
		} else if t1.Is(token.Symbol('^')) {
			return p.ParseCustomPointerType()
		}
	case token.Identifier:
		return p.ParseTypeForIdentifier()
	case token.NumeralInt, token.NumeralReal, token.CharacterString:
		return p.ParseConstSubrageType()
	case token.ReservedWord:
		switch t1.Value() {
		case "PACKED", "ARRAY", "SET", "RECORD", "FILE":
			return p.ParseStrucType()
		case "FUNCTION", "PROCEDURE":
			return p.ParseProcedureType()
		case "CLASS":
			return p.ParseClassType()
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

func (p *Parser) ParseTypeId() (*ast.TypeId, error) {
	if res, err := p.parseTypeIdWithUnit(); res != nil || err != nil {
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

func (p *Parser) ParseCustomPointerType() (*ast.CustomPointerType, error) {
	if _, err := p.Current(token.Symbol('^')); err != nil {
		return nil, err
	}
	p.NextToken()
	typ, err := p.ParseTypeId()
	if err != nil {
		return nil, err
	}
	return &ast.CustomPointerType{TypeId: typ}, nil
}
