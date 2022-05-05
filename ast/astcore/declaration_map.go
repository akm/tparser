package astcore

import (
	"strings"
)

type DeclarationMap interface {
	Get(name string) *Declaration
	Set(*Declaration)
	SetDecl(Decl)
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
