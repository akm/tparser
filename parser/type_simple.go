package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) ParseRealType(required bool) (ast.RealType, error) {
	t := p.CurrentToken()
	decl := ast.EmbeddedTypeDecl(ast.EtkReal, t.Value())
	if decl != nil {
		p.NextToken()
		return ast.NewTypeId(p.NewIdent(t), decl), nil
	} else if required {
		return nil, p.TokenErrorf("Unsupported token %s for RealType", t)
	} else {
		return nil, nil
	}
}

func (p *Parser) ParseOrdIdent(required bool) (ast.OrdIdent, error) {
	t := p.CurrentToken()
	decl := ast.EmbeddedTypeDecl(ast.EtkOrdIdent, t.Value())
	if decl != nil {
		p.NextToken()
		return ast.NewTypeId(p.NewIdent(t), decl), nil
	} else if required {
		return nil, p.TokenErrorf("Unsupported token %s for OrdIdent", t)
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
		qualId := ast.NewQualId(nil, p.NewIdentRef(t1))
		return &ast.SubrangeType{
			Low:  ast.NewConstExpr(qualId),
			High: expr,
		}, nil
	} else {
		defer rollback()
		if required {
			return nil, p.TokenErrorf("Unsupported token %s, %s for SubrangeType", t1, t2.RawString())
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
		Low:  lowExpr,
		High: highExpr,
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
			return nil, p.TokenErrorf("Unsupported token %s for EnumeratedType", t)
		}
	}
	p.NextToken()
	return res, nil
}

func (p *Parser) ParseEnumeratedTypeElement() (*ast.EnumeratedTypeElement, error) {
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.EnumeratedTypeElement{Ident: p.NewIdent(ident)}
	p.NextToken()
	if p.CurrentToken().Is(token.Symbol('=')) {
		p.NextToken()
		expr, err := p.ParseConstExpr()
		if err != nil {
			return nil, err
		}
		res.ConstExpr = expr
	}

	if err := p.context.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseTypeAsOrdinalType() (ast.OrdinalType, error) {
	t0 := p.CurrentToken()
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	if ordinalType, ok := typ.(ast.OrdinalType); !ok {
		return nil, errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
	} else if !ordinalType.IsOrdinalType() {
		return nil, errors.Errorf("Expected OrdinalType, got %T at %s", typ, p.PlaceString(t0))
	} else {
		return ordinalType, nil
	}
}
