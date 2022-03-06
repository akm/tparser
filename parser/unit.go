package parser

import (
	"github.com/akm/opparser/ast"
	"github.com/akm/opparser/token"
)

func (p *Parser) ParseUnit() (*ast.Unit, error) {
	// startToken := p.curr
	ident, err := p.get(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Unit{
		Ident: ast.Ident(ident.Text()),
	}
	t := p.next()
	if t.Type == token.PortabilityDirective {
		s := ast.PortabilityDirective(t.Text())
		res.PortabilityDirective = &s
		t = p.next()
	}
	if !t.Is(token.Symbol(';')) {
		return nil, err
	}
	if _, err := p.get(token.ReservedWord.HasText("interface")); err != nil {
		return nil, err
	}
	intf, err := p.ParseInterfaceSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.get(token.ReservedWord.HasText("implementation")); err != nil {
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
	return &ast.InterfaceSection{}, nil
}

func (p *Parser) ParseImplementationSection() (*ast.ImplementationSection, error) {
	return &ast.ImplementationSection{}, nil
}