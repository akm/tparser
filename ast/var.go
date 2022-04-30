package ast

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
	IdentList            IdentList
	Type                 Type
	Absolute             VarDeclAbsolute
	ConstExpr            *ConstExpr
	PortabilityDirective *PortabilityDirective
}

func (m *VarDecl) Children() Nodes {
	return Nodes{m.IdentList, m.Type, m.Absolute, m.ConstExpr}
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
	IdentList IdentList
	Type      Type
}

func (m *ThreadVarDecl) Children() Nodes {
	return Nodes{m.IdentList, m.Type}
}
