package ast

// - ConstSection
//   ```
//   CONST (ConstantDecl ';')...
//   ```
func (ConstSection) canBeInterfaceDecl() {}

type ConstSection []*ConstantDecl

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
}

func (m *ConstantDecl) Children() Nodes {
	r := Nodes{m.Ident}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	r = append(r, m.ConstExpr)
	return r
}

// - ConstExpr
//   ```
//   <constant-expression>
//   ```
type ConstExpr = Expression

func NewConstExpr(arg interface{}) *ConstExpr {
	return NewExpression(arg)
}
