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

	defer p.StackContext()()

	// startToken := p.curr
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Unit{
		Ident: ast.NewIdent(ident),
	}

	t := p.NextToken()
	if t.Is(token.PortabilityDirective) {
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

	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("INITIALIZATION")) {
		if initSection, err := p.ParseInitSection(); err != nil {
			return nil, err
		} else if initSection != nil {
			res.InitSection = initSection
		}
	}

	if _, err := p.Current(token.ReservedWord.HasKeyword("END")); err != nil {
		return nil, err
	}
	if _, err := p.Next(token.Symbol('.')); err != nil {
		return nil, err
	}
	res.InterfaceSection = intf
	res.ImplementationSection = impl
	p.context.SetDecl(res)
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
		res.UsesClause = usesClause
		p.context.AddUnitIdentifiers(usesClause.IdentList().Names()...)
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
			section, err := p.ParseTypeSection(true)
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "VAR":
			section, err := p.ParseVarSection(true)
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "THREADVAR":
			section, err := p.ParseThreadVarSection(true)
			if err != nil {
				return nil, err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "CONST":
			section, err := p.ParseConstSection(true)
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
	p.NextToken()

	defer p.StackContext()()

	res := &ast.ImplementationSection{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		usesClause, err := p.ParseUsesClause()
		if err != nil {
			return nil, err
		}
		res.UsesClause = usesClause
		p.context.AddUnitIdentifiers(usesClause.IdentList().Names()...)
		p.NextToken()
	}

	if declSections, err := p.ParseDeclSections(); err != nil {
		return nil, err
	} else if len(declSections) > 0 {
		res.DeclSections = declSections
	}

	if exportsStmt, err := p.ParseExportsStmts(); err != nil {
		return nil, err
	} else if exportsStmt != nil {
		res.ExportsStmts = exportsStmt
	}

	if p.CurrentToken().Is(token.Symbol(';')) {
		p.NextToken()
	}

	return res, nil
}

func (p *Parser) ParseInitSection() (*ast.InitSection, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("INITIALIZATION")); err != nil {
		return nil, err
	}
	p.NextToken()

	res := &ast.InitSection{}

	predFinalization := token.ReservedWord.HasKeyword("FINALIZATION")
	predEnd := token.ReservedWord.HasKeyword("END")
	terminator := token.Some(predFinalization, predEnd)

	if stmtList, err := p.ParseStmtList(terminator); err != nil {
		return nil, err
	} else if stmtList != nil && len(stmtList) > 0 {
		res.InitializationStmts = stmtList
	}

	if p.CurrentToken().Is(predFinalization) {
		p.NextToken()
		if stmtList, err := p.ParseStmtList(predEnd); err != nil {
			return nil, err
		} else if stmtList != nil && len(stmtList) > 0 {
			res.FinalizationStmts = stmtList
		}
	}

	return res, nil
}

func (p *Parser) ParseQualIds() (ast.QualIds, error) {
	res := ast.QualIds{}
	separator := token.Symbol(',')
	err := p.Until(token.Not(separator), separator, func() error {
		qualId, err := p.ParseQualId()
		if err != nil {
			return err
		}
		res = append(res, qualId)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, nil
	} else {
		return res, nil
	}
}

func (p *Parser) ParseQualId() (*ast.QualId, error) {
	if _, err := p.Current(token.Some(token.Identifier)); err != nil {
		return nil, err
	}
	name1 := p.CurrentToken()
	if p.context.IsUnitIdentifier(name1) {
		if _, err := p.Next(token.Symbol('.')); err != nil {
			return nil, err
		}
		name2, err := p.Next(token.Identifier)
		if err != nil {
			return nil, err
		}
		// TODO find Declaration from Unit in context
		p.NextToken()
		return ast.NewQualId(ast.NewUnitId(name1), ast.NewIdent(name2)), nil
	} else {
		p.NextToken()
		d := p.context.Get(name1.RawString())
		return ast.NewQualId(nil, ast.NewIdent(name1), d), nil
	}
}
