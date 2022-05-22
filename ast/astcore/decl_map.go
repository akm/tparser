package astcore

import (
	"strings"
)

type DeclMap interface {
	Get(name string) *Decl
	Set(DeclNode) error
}

type DeclMapImpl map[string]*Decl

func NewDeclarationMap() DeclMap {
	return make(DeclMapImpl)
}

func (m DeclMapImpl) Set(decl DeclNode) error {
	for _, i := range decl.ToDeclarations() {
		s := m.regularize(i.Ident.Name)
		// if s == "bar" && fmt.Sprintf("%T", i.Node) == "*ast.Unit" {
		// 	err := errors.Errorf("bar found")
		// 	fmt.Printf("%+v\n", err)
		// }
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
