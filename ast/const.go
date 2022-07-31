package ast

import "github.com/akm/tparser/ast/astcore"

// - ConstSection
//   ```
//   CONST (ConstantDecl ';')...
//   ```
type ConstSection []*ConstantDecl // must implement InterfaceDecl

var _ InterfaceDecl = (ConstSection)(nil)
var _ DeclSection = (ConstSection)(nil)

func (ConstSection) canBeInterfaceDecl() {}
func (ConstSection) canBeDeclSection()   {}
func (s ConstSection) Children() Nodes {
	r := make(Nodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}
func (s ConstSection) GetDeclNodes() astcore.DeclNodes {
	r := make(astcore.DeclNodes, len(s))
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
}

var _ astcore.DeclNode = (*ConstantDecl)(nil)

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

type ConstExprs = ExprList
