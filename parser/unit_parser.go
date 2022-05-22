package parser

import (
	"io/ioutil"
	"os"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/token"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type UnitParser struct {
	*Parser
	context *UnitContext
	Unit    *ast.Unit
}

func NewUnitParser(ctx *UnitContext) *UnitParser {
	return &UnitParser{Parser: NewParser(ctx), context: ctx}
}

func (p *UnitParser) LoadFile() error {
	fp, err := os.Open(p.context.Path)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %q", p.context.Path)
	}
	defer fp.Close()

	decoder := japanese.ShiftJIS.NewDecoder()
	str, err := ioutil.ReadAll(transform.NewReader(fp, decoder))
	if err != nil {
		return err
	}

	runes := []rune(string(str))
	p.SetText(&runes)
	p.Parser.NextToken()
	return nil
}

// ParseUnit method is not deleted for tests.
// Don't use this method not for test.
func (p *UnitParser) ParseUnit() (*ast.Unit, error) {
	res, err := p.ParseUnitIdentAndIntfUses()
	if err != nil {
		return nil, err
	}
	if err := p.ParseUnitIntfBody(); err != nil {
		return nil, err
	}

	if err := p.ParseImplUses(); err != nil {
		return nil, err
	}
	if err := p.ParseImplBody(); err != nil {
		return nil, err
	}

	if err := p.ParseUnitEnd(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *UnitParser) ProcessIdentAndIntfUses() error {
	_, err := p.ParseUnitIdentAndIntfUses()
	return err
}

func (p *UnitParser) ParseUnitIdentAndIntfUses() (*ast.Unit, error) {
	res, err := p.ParseUnitIdent()
	if err != nil {
		return nil, err
	}
	if err := p.ParseUnitIntfUses(); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *UnitParser) ParseUnitIdent() (*ast.Unit, error) {
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
	p.Unit = res

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

func (p *UnitParser) ParseUnitIntfUses() error {
	intf, err := p.ParseInterfaceSectionUses()
	if err != nil {
		return err
	}
	p.Unit.InterfaceSection = intf
	return nil
}

func (p *UnitParser) ParseInterfaceSectionUses() (*ast.InterfaceSection, error) {
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
		p.NextToken()
	}
	return res, nil
}

func (m *UnitParser) ProcessIntfBody() error {
	// Import decls after resovling unit load order and parsing units which this unit uses.
	if err := m.context.ImportUnitDecls(m.Unit.InterfaceSection.UsesClause); err != nil {
		return err
	}
	if err := m.context.Set(m.Unit); err != nil {
		return err
	}

	// Parse rest of interface Section (except USES clause)
	if err := m.ParseUnitIntfBody(); err != nil {
		return err
	}

	m.Unit.DeclarationMap = m.context.DeclMap
	return nil
}

func (p *UnitParser) ParseUnitIntfBody() error {
	if err := p.ParseInterfaceSectionDecls(); err != nil {
		return err
	}
	p.Unit.DeclarationMap = p.context.GetDeclarationMap()
	return nil
}

func (p *UnitParser) ParseUnitEnd() error {
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("INITIALIZATION")) {
		if initSection, err := p.ParseInitSection(); err != nil {
			return err
		} else if initSection != nil {
			p.Unit.InitSection = initSection
		}
	}

	if _, err := p.Current(token.ReservedWord.HasKeyword("END")); err != nil {
		return err
	}
	if _, err := p.Next(token.Symbol('.')); err != nil {
		return err
	}
	return nil
}

func (p *UnitParser) ParseInterfaceSectionDecls() error {
	res := p.Unit.InterfaceSection
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

func (m *UnitParser) ProcessImplAndInit() error {
	if err := m.ParseImplUses(); err != nil {
		return err
	}
	// defer m.Parser.StackContext()()

	if err := m.context.ImportUnitDecls(m.Unit.ImplementationSection.UsesClause); err != nil {
		return err
	}

	if err := m.ParseImplBody(); err != nil {
		return err
	}
	if err := m.ParseUnitEnd(); err != nil {
		return err
	}

	return nil
}

func (p *UnitParser) ParseImplUses() error {
	if _, err := p.Current(token.ReservedWord.HasKeyword("IMPLEMENTATION")); err != nil {
		return err
	}
	p.NextToken()

	impl := &ast.ImplementationSection{}
	if p.CurrentToken().Is(token.ReservedWord.HasKeyword("USES")) {
		usesClause, err := p.ParseUsesClause()
		if err != nil {
			return err
		}
		impl.UsesClause = usesClause
		p.NextToken()
	}
	p.Unit.ImplementationSection = impl
	return nil
}

func (p *UnitParser) ParseImplBody() error {
	if declSections, err := p.ParseDeclSections(); err != nil {
		return err
	} else if len(declSections) > 0 {
		p.Unit.ImplementationSection.DeclSections = declSections
	}

	if exportsStmt, err := p.ParseExportsStmts(); err != nil {
		return err
	} else if exportsStmt != nil {
		p.Unit.ImplementationSection.ExportsStmts = exportsStmt
	}

	if p.CurrentToken().Is(token.Symbol(';')) {
		p.NextToken()
	}

	return nil
}

func (p *UnitParser) ParseInitSection() (*ast.InitSection, error) {
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
