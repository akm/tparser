package ast

import "github.com/akm/tparser/ast/astcore"

// - VarSection
//   ```
//   VAR (VarDecl ';')...
//   ```
func (VarSection) canBeInterfaceDecl() {}

type VarSection []*VarDecl

func (s VarSection) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
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
	astcore.Decl
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

func (m *VarDecl) ToDeclarations() astcore.Declarations {
	return astcore.NewDeclarations(m.IdentList, m)
}

type VarDeclAbsolute interface {
	Node
	isVarDeclAbsolute()
}

func (VarDeclAbsoluteIdent) isVarDeclAbsolute() {}

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

func (*VarDeclAbsoluteIdent) Children() Nodes {
	return Nodes{}
}

func (*VarDeclAbsoluteConstExpr) isVarDeclAbsolute() {}

type VarDeclAbsoluteConstExpr ConstExpr

// threadvar X: Integer;
// Thread-variable declarations
// • cannot occur within a procedure or function.
// • cannot include initializations.
// • cannot specify the absolute directive.
func (ThreadVarSection) canBeInterfaceDecl() {}

type ThreadVarSection []*ThreadVarDecl

func (s ThreadVarSection) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

type ThreadVarDecl struct {
	IdentList
	Type Type
	astcore.Decl
}

func (m *ThreadVarDecl) Children() Nodes {
	r := Nodes{m.IdentList}
	if m.Type != nil {
		r = append(r, m.Type)
	}
	return r
}

func (m *ThreadVarDecl) ToDeclarations() astcore.Declarations {
	return astcore.NewDeclarations(m.IdentList, m)
}
