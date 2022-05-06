package ast

// - Program
//   ```
//   [PROGRAM Ident ['(' IdentList ')'] ';']
//   ProgramBlock '.'
//   ```
// In standard Pascal, a program heading can include parameters after the program name:
//     program Calc(input, output);
// Borland’s Object Pascal compiler ignores these parameters.
type Program struct {
	Node
	*Ident
	// IdentList    *IdentList // Borland’s Object Pascal compiler ignores these parameters.
	ProgramBlock *ProgramBlock
}

func (m *Program) Children() Nodes {
	res := Nodes{}
	if m.Ident != nil {
		res = append(res, m.Ident)
	}
	res = append(res, m.ProgramBlock)
	return res
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
