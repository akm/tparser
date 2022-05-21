package astcore

import (
	"strings"
)

type DeclarationMap interface {
	Get(name string) *Decl
	Set(DeclNode)
}

type declarationMapImpl map[string]*Decl

func NewDeclarationMap() DeclarationMap {
	return make(declarationMapImpl)
}

func (m declarationMapImpl) Set(decl DeclNode) {
	for _, i := range decl.ToDeclarations() {
		m[m.regularize(i.Ident.Name)] = i
	}
}

func (m declarationMapImpl) Get(name string) *Decl {
	return m[m.regularize(name)]
}

func (m declarationMapImpl) regularize(name string) string {
	return strings.ToLower(name)
}

type CompositeDeclarationMap struct {
	maps []DeclarationMap
}

func NewCompositeDeclarationMap(maps ...DeclarationMap) *CompositeDeclarationMap {
	return &CompositeDeclarationMap{maps: maps}
}

func (c *CompositeDeclarationMap) Get(name string) *Decl {
	for _, m := range c.maps {
		if d := m.Get(name); d != nil {
			return d
		}
	}
	return nil
}

func (c *CompositeDeclarationMap) Set(decl DeclNode) {
	c.maps[0].Set(decl)
}
