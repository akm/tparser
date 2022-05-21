package parser

import (
	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
)

func (p *Parser) IsUnitIdentifier() bool {
	return p.context.IsUnitIdentifier(p.CurrentToken())
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
		unitDecl := p.context.Get(name1.Value())
		if unitDecl == nil {
			return nil, p.TokenErrorf("undefined unit %s", name1)
		}
		if !IsUnitDeclaration(unitDecl) {
			return nil, p.TokenErrorf("%s is not a unit", name1)
		}
		unit := unitDecl.Node.(*ast.Unit)
		decl := unit.DeclarationMap.Get(name2.Value())
		if decl == nil {
			return nil, p.TokenErrorf("undefined identifier %s in unit %s", name2, name1.Value())
		}
		p.NextToken()
		return &ast.QualId{
			UnitId: &ast.IdentRef{Ident: ast.NewIdent(name1), Ref: unitDecl},
			Ident:  &ast.IdentRef{Ident: ast.NewIdent(name2), Ref: decl},
		}, nil
	} else {
		p.NextToken()
		return ast.NewQualId(nil, p.NewIdentRef(name1)), nil
	}
}

// ParseUnit method is not deleted for tests.
// Don't use this method not for test.
func (p *Parser) ParseUnit() (*ast.Unit, error) {
	res, err := p.ParseUnitIdentAndIntfUses()
	if err != nil {
		return nil, err
	}
	if err := p.ParseUnitIntfBody(res); err != nil {
		return nil, err
	}

	if err := p.ParseImplUses(res); err != nil {
		return nil, err
	}
	if err := p.ParseImplBody(res); err != nil {
		return nil, err
	}

	if err := p.ParseUnitEnd(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Parser) ParseUnitIdentAndIntfUses() (*ast.Unit, error) {
	res, err := p.ParseUnitIdent()
	if err != nil {
		return nil, err
	}
	if err := p.ParseUnitIntfUses(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) ParseUnitIdent() (*ast.Unit, error) {
	if _, err := p.Current(token.ReservedWord.HasKeyword("UNIT")); err != nil {
		return nil, err
	}

	// defer p.StackContext()()

	// startToken := p.curr
	ident, err := p.Next(token.Identifier)
	if err != nil {
		return nil, err
	}
	res := &ast.Unit{
		Path:  p.context.GetPath(),
		Ident: p.NewIdent(ident),
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
	return res, nil
}

func (p *Parser) ParseUnitIntfUses(res *ast.Unit) error {
	intf, err := p.ParseInterfaceSectionUses()
	if err != nil {
		return err
	}
	res.InterfaceSection = intf
	return nil
}

func (p *Parser) ParseUnitIntfBody(res *ast.Unit) error {
	if err := p.ParseInterfaceSectionDecls(res.InterfaceSection); err != nil {
		return err
	}
	res.DeclarationMap = p.context.GetDeclarationMap()
	return nil
}

func (p *Parser) ParseUnitEnd(res *ast.Unit) error {
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("INITIALIZATION")) {
		if initSection, err := p.ParseInitSection(); err != nil {
			return err
		} else if initSection != nil {
			res.InitSection = initSection
		}
	}

	if _, err := p.Current(token.ReservedWord.HasKeyword("END")); err != nil {
		return err
	}
	if _, err := p.Next(token.Symbol('.')); err != nil {
		return err
	}
	p.context.Set(res)
	return nil
}

func (p *Parser) ParseInterfaceSectionUses() (*ast.InterfaceSection, error) {
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
	return res, nil
}

func (p *Parser) ParseInterfaceSectionDecls(res *ast.InterfaceSection) error {
	res.InterfaceDecls = []ast.InterfaceDecl{}
	defer func() {
		if len(res.InterfaceDecls) == 0 {
			res.InterfaceDecls = nil
		}
	}()

	for {
		t := p.CurrentToken()
		if t.Is(token.EOF) {
			return nil
		}
		if !t.Is(token.ReservedWord) {
			return p.TokenErrorf("expects reserved word but got %s", t)
		}
		switch t.Value() {
		case "TYPE":
			section, err := p.ParseTypeSection(true)
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "VAR":
			section, err := p.ParseVarSection(true)
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "THREADVAR":
			section, err := p.ParseThreadVarSection(true)
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "CONST":
			section, err := p.ParseConstSection(true)
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "FUNCTION":
			section, err := p.ParseExportedHeading()
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		case "PROCEDURE":
			section, err := p.ParseExportedHeading()
			if err != nil {
				return err
			}
			res.InterfaceDecls = append(res.InterfaceDecls, section)
			continue
		}
		break
	}
	return nil
}

func (p *Parser) ParseImplUses(res *ast.Unit) error {
	if _, err := p.Current(token.ReservedWord.HasKeyword("IMPLEMENTATION")); err != nil {
		return err
	}
	p.NextToken()

	defer p.StackContext()()

	impl := &ast.ImplementationSection{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		usesClause, err := p.ParseUsesClause()
		if err != nil {
			return err
		}
		impl.UsesClause = usesClause
		p.context.AddUnitIdentifiers(usesClause.IdentList().Names()...)
		p.NextToken()
	}
	res.ImplementationSection = impl
	return nil
}

func (p *Parser) ParseImplBody(res *ast.Unit) error {
	if declSections, err := p.ParseDeclSections(); err != nil {
		return err
	} else if len(declSections) > 0 {
		res.ImplementationSection.DeclSections = declSections
	}

	if exportsStmt, err := p.ParseExportsStmts(); err != nil {
		return err
	} else if exportsStmt != nil {
		res.ImplementationSection.ExportsStmts = exportsStmt
	}

	if p.CurrentToken().Is(token.Symbol(';')) {
		p.NextToken()
	}

	return nil
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

type UnitParser struct {
	*Parser
	context *UnitContext
}

func NewUnitParser(text *[]rune, ctx *UnitContext) *UnitParser {
	return &UnitParser{
		Parser:  NewParser(text, ctx),
		context: ctx,
	}
}
