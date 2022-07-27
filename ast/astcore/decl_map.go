package astcore

import (
	"strings"

	"github.com/pkg/errors"
)

type DeclMap interface {
	Get(name string) *Decl
	Set(DeclNode) error
	Overwrite(name string, decl *Decl)
}

type DeclMapImpl map[string]*Decl

var _ DeclMap = (DeclMapImpl)(nil)

func NewDeclMap() DeclMap {
	return make(DeclMapImpl)
}

func (m DeclMapImpl) Overwrite(name string, decl *Decl) {
	m[m.regularize(name)] = decl
}

func (m DeclMapImpl) Set(decl DeclNode) error {
	for _, i := range decl.ToDeclarations() {
		s := m.regularize(i.Ident.Name)
		// if s == "bar" && fmt.Sprintf("%T", i.Node) == "*ast.Unit" {
		// 	err := errors.Errorf("bar found")
		// 	fmt.Printf("%+v\n", err)
		// }
		if old, ok := m[s]; ok {
			return errors.Errorf("duplicate declaration: %s at %s but was %s", i.Ident.Name, i.Ident.Location.String(), old.Location.String())
		}
		m[s] = i
	}
	return nil
}

func (m DeclMapImpl) Get(name string) *Decl {
	return m[m.regularize(name)]
}

func (m DeclMapImpl) regularize(name string) string {
	return strings.ToLower(name)
}

type ChainedDeclMap struct {
	Parent DeclMap
	Impl   DeclMapImpl
}

var _ DeclMap = (*ChainedDeclMap)(nil)

func NewChainedDeclMap(parent DeclMap) *ChainedDeclMap {
	return &ChainedDeclMap{Parent: parent, Impl: DeclMapImpl{}}
}

func (m *ChainedDeclMap) Get(name string) *Decl {
	if d := m.Impl.Get(name); d != nil {
		return d
	}
	return m.Parent.Get(name)
}

func (m *ChainedDeclMap) Set(n DeclNode) error {
	return m.Impl.Set(n)
}

func (m *ChainedDeclMap) Overwrite(name string, decl *Decl) {
	m.Impl.Overwrite(name, decl)
}

type CompositeDeclMap struct {
	maps []DeclMap
}

var _ DeclMap = (*CompositeDeclMap)(nil)

func NewCompositeDeclMap(maps ...DeclMap) *CompositeDeclMap {
	return &CompositeDeclMap{maps: maps}
}

func (c *CompositeDeclMap) Get(name string) *Decl {
	for _, m := range c.maps {
		if d := m.Get(name); d != nil {
			return d
		}
	}
	return nil
}

func (c *CompositeDeclMap) Set(decl DeclNode) error {
	return c.maps[0].Set(decl)
}

func (c *CompositeDeclMap) Overwrite(name string, decl *Decl) {
	c.maps[0].Overwrite(name, decl)
}

type DeclMaps []DeclMap

func (s DeclMaps) Reverse() DeclMaps {
	r := make(DeclMaps, len(s))
	for i, m := range s {
		r[len(s)-i-1] = m
	}
	return r
}
