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
	astcore.DeclNode
	*FunctionHeading
	Directives           []Directive
	ExternalOptions      *ExternalOptions
	PortabilityDirective *PortabilityDirective
	Block                *Block
}

func (*FunctionDecl) canBeDeclSection()       {}
func (*FunctionDecl) isProcedureDeclSection() {}
func (m *FunctionDecl) Children() Nodes {
	return Nodes{m.FunctionHeading, m.Block}
}
func (m *FunctionDecl) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}
