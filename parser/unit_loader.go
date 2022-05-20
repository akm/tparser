package parser

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	toposort "github.com/philopon/go-toposort"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type UnitLoader struct {
	ctx *UnitContext
	*Parser
	Unit *ast.Unit
}

func NewUnitLoader(ctx *UnitContext) *UnitLoader {
	return &UnitLoader{ctx: ctx}
}

func (m *UnitLoader) LoadFile() error {
	fp, err := os.Open(m.ctx.Path)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %q", m.ctx.Path)
	}
	defer fp.Close()

	decoder := japanese.ShiftJIS.NewDecoder()
	str, err := ioutil.ReadAll(transform.NewReader(fp, decoder))
	if err != nil {
		return err
	}

	runes := []rune(string(str))
	m.Parser = NewParser(&runes, m.ctx)
	return nil
}

func (m *UnitLoader) ProcessIdentAndIntfUses() error {
	m.Parser.NextToken()
	if u, err := m.Parser.ParseUnitIdentAndIntfUses(); err != nil {
		return err
	} else {
		m.Unit = u
		return nil
	}
}

func (m *UnitLoader) LoadBody() error {
	units := ast.Units{}
	parentUnits := m.ctx.Parent.Units
	for _, unitRef := range m.Unit.InterfaceSection.UsesClause {
		if u := parentUnits.ByName(unitRef.Ident.Name); u != nil {
			units = append(units, u)
		}
	}
	localMap := astcore.NewDeclarationMap()
	localMap.SetDecl(m.Unit)
	maps := []astcore.DeclarationMap{localMap}
	for _, unit := range units {
		localMap.SetDecl(unit)
		// TODO declMapに追加する順番はこれでOK？
		// 無関係のユニットAとBに、同じ名前の型や変数が定義されていて、USES A, B; となっていた場合
		// コンテキスト上ではどちらが有効になるのかを確認する
		maps = append(maps, unit.DeclarationMap)
	}
	m.ctx.DeclarationMap = astcore.NewCompositeDeclarationMap(maps...)

	// Parse rest of interface Section (except USES clause)
	if err := m.Parser.ParseUnitBody(m.Unit); err != nil {
		return err
	}

	m.Unit.DeclarationMap = m.ctx.DeclarationMap
	return nil
}

func (m *UnitLoader) LoadTail() error {
	return m.Parser.ParseUnitTail(m.Unit)
}

type UnitLoaders []*UnitLoader

func (m UnitLoaders) UnitNames() ext.StringSet {
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

func (m UnitLoaders) LoaderMap() map[string]*UnitLoader {
	r := map[string]*UnitLoader{}
	for _, loader := range m {
		name := strings.ToLower(loader.Unit.Name)
		r[name] = loader
	}
	return r
}

func (m UnitLoaders) Graph() *toposort.Graph {
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

func (m UnitLoaders) Sort() (UnitLoaders, error) {
	loaderMap := m.LoaderMap()

	graph := m.Graph()
	order, ok := graph.Toposort()
	if !ok {
		return nil, errors.Errorf("cyclic dependency detected")
	}

	r := make(UnitLoaders, 0, len(m))
	for _, name := range order {
		if loader, ok := loaderMap[name]; ok {
			r = append(r, loader)
		}
	}
	return r, nil
}

func (m UnitLoaders) DeclarationMaps() []astcore.DeclarationMap {
	r := make([]astcore.DeclarationMap, len(m))
	for i, loader := range m {
		r[i] = loader.Unit.DeclarationMap
	}
	return r
}
