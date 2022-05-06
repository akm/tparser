package ast

// - Block
//   ```
//   [DeclSection]
//   [ExportsStmt]...
//   CompoundStmt
//   [ExportsStmt]...
//   ```
type Block struct {
	Node
	DeclSections  DeclSections
	ExportsStmts1 ExportsStmts
	CompoundStmt  *CompoundStmt
	ExportsStmts2 ExportsStmts
}

func (m *Block) Children() Nodes {
	return Nodes{m.DeclSections, m.ExportsStmts1, m.CompoundStmt, m.ExportsStmts2}
}

type ExportsStmts []*ExportsStmt // must implements Node
func (s ExportsStmts) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// - ExportsStmt
//   ```
//   EXPORTS ExportsItem [, ExportsItem]...
//   ```
type ExportsStmt struct {
	Node
	ExportsItems []*ExportsItem
}

func (m *ExportsStmt) Children() Nodes {
	r := make(Nodes, len(m.ExportsItems))
	for idx, i := range m.ExportsItems {
		r[idx] = i
	}
	return r
}

// - ExportsItem
//   ```
//   Ident [NAME|INDEX “‘” ConstExpr “‘”]
//         [INDEX|NAME “‘” ConstExpr “‘”]
//   ```
type ExportsItem struct {
	Node
	*Ident
	Name  *ConstExpr
	Index *ConstExpr
}

func (m *ExportsItem) Children() Nodes {
	return Nodes{m.Ident, m.Name, m.Index}
}

// - DeclSection
//   ```
//   LabelDeclSection
//   ```
//   ```
//   ConstSection
//   ```
//   ```
//   TypeSection
//   ```
//   ```
//   VarSection
//   ```
//   ```
//   ProcedureDeclSection
//   ```
type DeclSection interface {
	Node
	canBeDeclSection()
}

type DeclSections []DeclSection // must implement Node

func (m DeclSections) Children() Nodes {
	r := make(Nodes, len(m))
	for idx, i := range m {
		r[idx] = i
	}
	return r
}

// - LabelDeclSection
//   ```
//   LABEL LabelId ';'
//   ```
type LabelDeclSection struct {
	*LabelId
}

func (*LabelDeclSection) canBeDeclSection() {}

type LabelId = Ident
