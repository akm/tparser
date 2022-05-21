package ast

import "github.com/akm/tparser/ast/astcore"

// - Program
//   ```
//   [PROGRAM Ident ['(' IdentList ')'] ';']
//   ProgramBlock '.'
//   ```
// In standard Pascal, a program heading can include parameters after the program name:
//     program Calc(input, output);
// Borland’s Object Pascal compiler ignores these parameters.
type Program struct {
	Path string
	*Ident
	// IdentList    *IdentList // Borland’s Object Pascal compiler ignores these parameters.
	ProgramBlock *ProgramBlock
	Goal
}

func (*Program) isGoal() {}
func (m *Program) GetPath() string {
	return m.Path
}
func (m *Program) Children() Nodes {
	res := Nodes{}
	if m.Ident != nil {
		res = append(res, m.Ident)
	}
	res = append(res, m.ProgramBlock)
	return res
}
func (m *Program) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.Ident, m)}
}

// - ProgramBlock
//   ```
//   [UsesClause]
//   Block
//   ```
type ProgramBlock struct {
	UsesClause UsesClause
	*Block
}

func (m *ProgramBlock) Children() Nodes {
	res := Nodes{}
	if m.UsesClause != nil {
		res = append(res, m.UsesClause)
	}
	res = append(res, m.Block)
	return res
}
