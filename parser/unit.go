package parser

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/akm/tparser/token"
	toposort "github.com/philopon/go-toposort"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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

func (m *UnitParser) ProcessIntfBody() error {
	units := ast.Units{}
	parentUnits := m.context.Parent.Units
	for _, unitRef := range m.Unit.InterfaceSection.UsesClause {
		if u := parentUnits.ByName(unitRef.Ident.Name); u != nil {
			units = append(units, u)
		}
	}
	localMap := astcore.NewDeclarationMap()
	localMap.Set(m.Unit)
	maps := []astcore.DeclMap{localMap}
	for _, unit := range units {
		localMap.Set(unit)
		// TODO declMapに追加する順番はこれでOK？
		// 無関係のユニットAとBに、同じ名前の型や変数が定義されていて、USES A, B; となっていた場合
		// コンテキスト上ではどちらが有効になるのかを確認する
		maps = append(maps, unit.DeclarationMap)
	}
	m.context.DeclMap = astcore.NewCompositeDeclarationMap(maps...)

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
	p.context.Set(p.Unit)
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
		p.context.AddUnitIdentifiers(usesClause.IdentList().Names()...)
		p.NextToken()
	}
	return res, nil
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

	parentUnits := m.context.Parent.Units
	localMap := astcore.NewDeclarationMap()
	localMap.Set(m.Unit)
	maps := []astcore.DeclMap{localMap}
	for _, unitRef := range m.Unit.ImplementationSection.UsesClause {
		if unit := parentUnits.ByName(unitRef.Ident.Name); unit != nil {
			localMap.Set(unit)
			maps = append(maps, unit.DeclarationMap)
		}
	}
	m.context.DeclMap = astcore.NewCompositeDeclarationMap(maps...)

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

type UnitParsers []*UnitParser

func (m UnitParsers) Units() ast.Units {
	r := make(ast.Units, len(m))
	for i, loader := range m {
		r[i] = loader.Unit
	}
	return r
}

func (m UnitParsers) UnitNames() ext.StringSet {
	unitNames := ext.Strings{}
	for _, loader := range m {
		name := strings.ToLower(loader.Unit.Name)
		unitNames = append(unitNames, name)
		for _, ref := range loader.Unit.InterfaceSection.UsesClause {
			unitNames = append(unitNames, strings.ToLower(ref.Name))
		}
	}
	return unitNames.Set()
}

func (m UnitParsers) Map() map[string]*UnitParser {
	r := map[string]*UnitParser{}
	for _, loader := range m {
		name := strings.ToLower(loader.Unit.Name)
		r[name] = loader
	}
	return r
}

func (m UnitParsers) Graph() *toposort.Graph {
	unitNames := m.UnitNames().Slice()
	graph := toposort.NewGraph(len(unitNames))
	graph.AddNodes(unitNames...)

	for _, loader := range m {
		for _, unitRef := range loader.Unit.InterfaceSection.UsesClause {
			graph.AddEdge(strings.ToLower(unitRef.Name), strings.ToLower(loader.Unit.Ident.Name))
		}
	}
	return graph
}

func (m UnitParsers) Sort() (UnitParsers, error) {
	loaderMap := m.Map()

	graph := m.Graph()
	order, ok := graph.Toposort()
	if !ok {
		return nil, errors.Errorf("cyclic dependency detected")
	}

	r := make(UnitParsers, 0, len(m))
	for _, name := range order {
		if loader, ok := loaderMap[name]; ok {
			r = append(r, loader)
		}
	}
	return r, nil
}

func (m UnitParsers) DeclarationMaps() []astcore.DeclMap {
	r := make([]astcore.DeclMap, len(m))
	for i, loader := range m {
		r[i] = loader.Unit.DeclarationMap
	}
	return r
}
