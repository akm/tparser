package ast

import "github.com/akm/tparser/ast/astcore"

// - VarSection
//   ```
//   VAR (VarDecl ';')...
//   ```
type VarSection []*VarDecl // must implement InterfaceDecl

func (VarSection) canBeInterfaceDecl() {}
func (VarSection) canBeDeclSection()   {}
func (s VarSection) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}
func (s VarSection) GetDeclNodes() astcore.DeclNodes {
	r := make(astcore.DeclNodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

// - VarDecl
//   - (On Windows)
//     ```
//     IdentList ':' Type [(ABSOLUTE (Ident | ConstExpr)) | '=' ConstExpr] [PortabilityDirective]
//     ```
//   - On Linux
//     ```
//     IdentList ':' Type [ABSOLUTE (Ident) | '=' ConstExpr] [PortabilityDirective]
//     ```
type VarDecl struct {
	IdentList
	Type                 Type
	Absolute             VarDeclAbsolute
	ConstExpr            *ConstExpr
	PortabilityDirective *PortabilityDirective
	astcore.DeclNode
}

func (m *VarDecl) Children() Nodes {
	r := Nodes{m.IdentList}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	if m.Absolute != nil {
		r = append(r, m.Absolute)
	}
	if m.ConstExpr != nil {
		r = append(r, m.ConstExpr)
	}
	return r
}

func (m *VarDecl) ToDeclarations() astcore.Decls {
	return astcore.NewDeclarations(m.IdentList, m)
}

type VarDeclAbsolute interface {
	Node
	isVarDeclAbsolute()
}

type VarDeclAbsoluteIdent Ident

func NewVarDeclAbsoluteIdent(arg interface{}) *VarDeclAbsoluteIdent {
	switch v := arg.(type) {
	case Ident:
		r := VarDeclAbsoluteIdent(v)
		return &r
	case *Ident:
		r := VarDeclAbsoluteIdent(*v)
		return &r
	default:
		return NewVarDeclAbsoluteIdent(NewIdentFrom(arg))
	}
}

func (VarDeclAbsoluteIdent) isVarDeclAbsolute() {}
func (*VarDeclAbsoluteIdent) Children() Nodes {
	return Nodes{}
}

type VarDeclAbsoluteConstExpr ConstExpr

func (*VarDeclAbsoluteConstExpr) isVarDeclAbsolute() {}

// threadvar X: Integer;
// Thread-variable declarations
// • cannot occur within a procedure or function.
// • cannot include initializations.
// • cannot specify the absolute directive.
type ThreadVarSection []*ThreadVarDecl // must implement InterfaceDecl

func (ThreadVarSection) canBeInterfaceDecl() {}
func (ThreadVarSection) canBeDeclSection()   {}
func (s ThreadVarSection) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}
func (s ThreadVarSection) GetDeclNodes() astcore.DeclNodes {
	r := make(astcore.DeclNodes, len(s))
	for i, m := range s {
		r[i] = m
	}
	return r
}

type ThreadVarDecl struct {
	IdentList
	Type Type
	astcore.DeclNode
}

func (m *ThreadVarDecl) Children() Nodes {
	r := Nodes{m.IdentList}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	return r
}

func (m *ThreadVarDecl) ToDeclarations() astcore.Decls {
	return astcore.NewDeclarations(m.IdentList, m)
}
