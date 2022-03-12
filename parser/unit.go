package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) ParseUnit() (*ast.Unit, error) {
	// startToken := p.curr
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Unit{
		Ident: ast.Ident(ident.Value()),
	}
	t := p.NextToken()
	if t.Type == token.PortabilityDirective {
		s := ast.PortabilityDirective(t.Value())
		res.PortabilityDirective = &s
		t = p.NextToken()
	}
	if !t.Is(token.Symbol(';')) {
		return nil, err
	}
	if _, err := p.Next(token.ReservedWord.HasKeyword("INTERFACE")); err != nil {
		return nil, err
	}
	intf, err := p.ParseInterfaceSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.Next(token.ReservedWord.HasKeyword("IMPLEMENTATION")); err != nil {
		return nil, err
	}
	impl, err := p.ParseImplementationSection()
	if err != nil {
		return nil, err
	}
	res.InterfaceSection = intf
	res.ImplementationSection = impl
	return res, nil
}

func (p *Parser) ParseInterfaceSection() (*ast.InterfaceSection, error) {
	res := &ast.InterfaceSection{}
	t := p.NextToken()
	if t.Is(token.ReservedWord.HasKeyword("USES")) {
		usesClause, err := p.ParseIdentClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = usesClause
		t = p.NextToken()
	}

	res.InterfaceDecls = []ast.InterfaceDecl{}
	defer func() {
		if len(res.InterfaceDecls) == 0 {
			res.InterfaceDecls = nil
		}
	}()

	for {
		t := p.CurrentToken()
		if t.Is(token.EOF) {
			return res, nil
		}
		if !t.Is(token.ReservedWord) {
			return nil, errors.Errorf("expects reserved word but got %s", t.String())
		}
		switch t.Value() {
		case "TYPE":
			section, err := p.ParseTypeSection()
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		}
		break
	}
	return res, nil
}

func (p *Parser) ParseImplementationSection() (*ast.ImplementationSection, error) {
	return &ast.ImplementationSection{}, nil
}
