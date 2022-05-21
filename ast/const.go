package ast

import "github.com/akm/tparser/ast/astcore"

// - ConstSection
//   ```
//   CONST (ConstantDecl ';')...
//   ```
type ConstSection []*ConstantDecl // must implement InterfaceDecl

func (ConstSection) canBeInterfaceDecl() {}
func (ConstSection) canBeDeclSection()   {}
func (s ConstSection) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - ConstantDecl
//   ```
//   Ident '=' ConstExpr [PortabilityDirective]
//   ```
//   ```
//   Ident ':' TypeId '=' TypedConstant [PortabilityDirective]
//   ```
type ConstantDecl struct {
	*Ident
	Type                 Type
	ConstExpr            *ConstExpr
	PortabilityDirective *PortabilityDirective
	astcore.DeclNode
}

func (m *ConstantDecl) Children() Nodes {
	r := Nodes{m.Ident}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	r = append(r, m.ConstExpr)
	return r
}

func (m *ConstantDecl) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

// - ConstExpr
//   ```
//   <constant-expression>
//   ```
type ConstExpr = Expression

func NewConstExpr(arg interface{}) *ConstExpr {
	return NewExpression(arg)
}
