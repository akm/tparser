package parser

import (
	"strings"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ext"
	"github.com/philopon/go-toposort"
	"github.com/pkg/errors"
)

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

func (m UnitParsers) DeclMaps() []astcore.DeclMap {
	r := make([]astcore.DeclMap, len(m))
	for i, loader := range m {
		r[i] = loader.Unit.DeclMap
	}
	return r
}
