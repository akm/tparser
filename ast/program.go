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
	*Ident
	// IdentList    *IdentList // Borland’s Object Pascal compiler ignores these parameters.
	ProgramBlock *ProgramBlock
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
