package ast

func (VarSection) canBeInterfaceDecl() {}

type VarSection []*VarDecl

func (s VarSection) Children() []Node {
	r := make([]Node, len(s))
	for i, v := range s {
		r[i] = v
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
	CodeBlockNode
	IdentList            IdentList
	Type                 Type
	Absolute             VarDeclAbsolute
	ConstExpr            *ConstExpr
	PortabilityDirective *PortabilityDirective
}

type VarDeclAbsolute interface {
	isVarDeclAbsolute()
}

func (VarDeclAbsoluteIdent) isVarDeclAbsolute() {}

type VarDeclAbsoluteIdent Ident

func NewVarDeclAbsoluteIdent(ident Ident) *VarDeclAbsoluteIdent {
	r := VarDeclAbsoluteIdent(ident)
	return &r
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

func (s ThreadVarSection) Children() []Node {
	r := make([]Node, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

type ThreadVarDecl struct {
	CodeBlockNode
	IdentList IdentList
	Type      Type
}
