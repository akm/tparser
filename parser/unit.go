package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
)

func (p *Parser) IsUnitIdentifier() bool {
	return p.context.IsUnitIdentifier(p.CurrentToken())
}

func (p *Parser) ParseUnit() (*ast.Unit, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("UNIT")); err != nil {
		return nil, err
	}
	// startToken := p.curr
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Unit{
		Ident: *ast.NewIdent(ident.Value()),
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
	p.NextToken()
	intf, err := p.ParseInterfaceSection()
	if err != nil {
		return nil, err
	}
	impl, err := p.ParseImplementationSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.Current(token.ReservedWord.HasKeyword("END")); err != nil {
		return nil, err
	}
	if _, err := p.Next(token.Symbol('.')); err != nil {
		return nil, err
	}
	res.InterfaceSection = intf
	res.ImplementationSection = impl
	return res, nil
}

func (p *Parser) ParseInterfaceSection() (*ast.InterfaceSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("INTERFACE")); err != nil {
		return nil, err
	}
	res := &ast.InterfaceSection{}
	t := p.NextToken()
	if t.Is(token.ReservedWord.HasKeyword("USES")) {
		usesClause, err := p.ParseUsesClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = &usesClause
		p.context.unitIdentifiers = append(p.context.unitIdentifiers, usesClause.IdentList().Names()...)
		p.NextToken()
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
		case "VAR":
			section, err := p.ParseVarSection()
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "THREADVAR":
			section, err := p.ParseThreadVarSection()
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "CONST":
			section, err := p.ParseConstSection()
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "FUNCTION":
			section, err := p.ParseExportedHeading()
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "PROCEDURE":
			section, err := p.ParseExportedHeading()
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
	if _, err := p.Current(token.ReservedWord.HasKeyword("IMPLEMENTATION")); err != nil {
		return nil, err
	}
	res := &ast.ImplementationSection{}
	for {
		t := p.NextToken()
		if t.Is(token.EOF) {
			return nil, errors.Errorf("expects END but got %s", t.String())
		}
		if t.Is(token.ReservedWord) {
			break
		}
	}
	return res, nil
}

func (p *Parser) ParseQualId() (*ast.QualId, error) {
	if _, err := p.Current(token.Some(token.Identifier, token.Directive)); err != nil {
		return nil, err
	}
	name1 := string(p.CurrentToken().Raw())
	p.NextToken()
	if p.CurrentToken().Is(token.Symbol('.')) {
		p.NextToken()
		if _, err := p.Current(token.Identifier); err != nil {
			return nil, err
		}
		name2 := p.CurrentToken().Raw()
		p.NextToken()
		return &ast.QualId{
			UnitId: ast.NewUnitId(name1),
			Ident:  *ast.NewIdent(name2),
		}, nil
	} else {
		return &ast.QualId{
			Ident: *ast.NewIdent(name1),
		}, nil
	}
}
