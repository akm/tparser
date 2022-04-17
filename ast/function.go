package ast

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
func (*FunctionDecl) canBeDeclSection()       {}
func (*FunctionDecl) isProcedureDeclSection() {}

type FunctionDecl struct {
	Heading              *FunctionHeading
	Directive            *Directive
	PortabilityDirective *PortabilityDirective
	Block                *Block
}
