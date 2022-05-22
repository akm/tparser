package astcore

import (
	"strings"
)

type DeclMap interface {
	Get(name string) *Decl
	Set(DeclNode)
}

type declMapImpl map[string]*Decl

func NewDeclarationMap() DeclMap {
	return make(declMapImpl)
}

func (m declMapImpl) Set(decl DeclNode) {
	for _, i := range decl.ToDeclarations() {
		m[m.regularize(i.Ident.Name)] = i
	}
}

func (m declMapImpl) Get(name string) *Decl {
	return m[m.regularize(name)]
}

func (m declMapImpl) regularize(name string) string {
	return strings.ToLower(name)
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

func (c *CompositeDeclMap) Set(decl DeclNode) {
	c.maps[0].Set(decl)
}
