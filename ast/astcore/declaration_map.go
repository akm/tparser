package astcore

import (
	"strings"

	"github.com/akm/tparser/ext"
)

type DeclarationMap interface {
	Get(name string) *Declaration
	Set(*Declaration)
	SetDecl(Decl)
	Keys() ext.Strings
}

type declarationMapImpl map[string]*Declaration

func NewDeclarationMap() DeclarationMap {
	return make(declarationMapImpl)
}

func (m declarationMapImpl) Set(d *Declaration) {
	m[m.regularize(d.Ident.Name)] = d
}

func (m declarationMapImpl) SetDecl(decl Decl) {
	for _, i := range decl.ToDeclarations() {
		m.Set(i)
	}
}

func (m declarationMapImpl) Get(name string) *Declaration {
	return m[m.regularize(name)]
}

func (m declarationMapImpl) regularize(name string) string {
	return strings.ToLower(name)
}

func (m declarationMapImpl) Keys() ext.Strings {
	keys := ext.Strings{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

type CompositeDeclarationMap struct {
	maps []DeclarationMap
}

func NewCompositeDeclarationMap(maps ...DeclarationMap) *CompositeDeclarationMap {
	return &CompositeDeclarationMap{maps: maps}
}

func (c *CompositeDeclarationMap) Get(name string) *Declaration {
	for _, m := range c.maps {
		if d := m.Get(name); d != nil {
			return d
		}
	}
	return nil
}

func (c *CompositeDeclarationMap) Set(d *Declaration) {
	c.maps[0].Set(d)
}

func (c *CompositeDeclarationMap) SetDecl(decl Decl) {
	c.maps[0].SetDecl(decl)
}

func (c *CompositeDeclarationMap) Keys() ext.Strings {
	r := ext.Strings{}
	for _, m := range c.maps {
		r = append(r, m.Keys()...)
	}
	return r
}
