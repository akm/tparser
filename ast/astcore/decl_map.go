package astcore

import (
	"strings"

	"github.com/pkg/errors"
)

type DeclMap interface {
	Get(name string) *Decl
	Set(DeclNode) error
}

type DeclMapImpl map[string]*Decl

func NewDeclMap() DeclMap {
	return make(DeclMapImpl)
}

func (m DeclMapImpl) Set(decl DeclNode) error {
	for _, i := range decl.ToDeclarations() {
		s := m.regularize(i.Ident.Name)
		// if s == "bar" && fmt.Sprintf("%T", i.Node) == "*ast.Unit" {
		// 	err := errors.Errorf("bar found")
		// 	fmt.Printf("%+v\n", err)
		// }
		if _, ok := m[s]; ok {
			return errors.Errorf("duplicate declaration: %s", i.Ident.Name)
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

type CompositeDeclMap struct {
	maps []DeclMap
}

func NewCompositeDeclarationMap(maps ...DeclMap) *CompositeDeclMap {
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
