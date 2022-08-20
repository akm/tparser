package ast

import "github.com/akm/tparser/ast/astcore"

// - ProcedureDeclSection
//   ```
//   ProcedureDecl
//   ```
//   ```
//   FunctionDecl
//   ```
type ProcedureDeclSection interface {
	DeclSection
	isProcedureDeclSection()
}

// - ProcedureDecl
//   ```
//   ProcedureHeading ';' [Directive] [PortabilityDirective]
//   Block ';'
//   ```
// - FunctionDecl
//   ```
//   FunctionHeading ';' [Directive] [PortabilityDirective]
//   Block ';'
//   ```
type FunctionDecl struct {
	*FunctionHeading
	Directives           []Directive // TODO : remove this field
	ExternalOptions      *ExternalOptions
	PortabilityDirective *PortabilityDirective
	Block                *Block
}

var _ astcore.DeclNode = (*FunctionDecl)(nil)
var _ DeclSection = (*FunctionDecl)(nil)
var _ ProcedureDeclSection = (*FunctionDecl)(nil)

func (*FunctionDecl) canBeDeclSection()       {}
func (*FunctionDecl) isProcedureDeclSection() {}
func (m *FunctionDecl) Children() Nodes {
	return Nodes{m.FunctionHeading, m.Block}
}
func (m *FunctionDecl) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}
