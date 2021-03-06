package ast

import "github.com/akm/tparser/ast/astcore"

// - Block
//   ```
//   [DeclSection]
//   [ExportsStmt]...
//   BlockBody
//   [ExportsStmt]...
//   ```
// - BlockBody
//   ```
//   CompoundStmt
//   ```
//   ```
//   AssemberStatement
//   ```
type Block struct {
	DeclSections  DeclSections
	ExportsStmts1 ExportsStmts
	Body          BlockBody
	ExportsStmts2 ExportsStmts
}

var _ Node = (*Block)(nil)

func (m *Block) Children() Nodes {
	res := Nodes{}
	if m.DeclSections != nil {
		res = append(res, m.DeclSections)
	}
	if m.ExportsStmts1 != nil {
		res = append(res, m.ExportsStmts1)
	}
	res = append(res, m.Body)
	if m.ExportsStmts2 != nil {
		res = append(res, m.ExportsStmts2)
	}
	return res
}

type ExportsStmts []*ExportsStmt

var _ Node = (ExportsStmts)(nil)

func (s ExportsStmts) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// BlockBody is CompoundStmt or AssemblerStatement
type BlockBody interface {
	StructStmt // extends StructsStmt
	isBlockBody()
}

// - ExportsStmt
//   ```
//   EXPORTS ExportsItem [, ExportsItem]...
//   ```
type ExportsStmt struct {
	ExportsItems []*ExportsItem
}

var _ Node = (*ExportsStmt)(nil)

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
	*Ident
	Name  *ConstExpr
	Index *ConstExpr
}

var _ Node = (*ExportsItem)(nil)

func (m *ExportsItem) Children() Nodes {
	res := Nodes{m.Ident}
	if m.Name != nil {
		res = append(res, m.Name)
	}
	if m.Index != nil {
		res = append(res, m.Index)
	}
	return res
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

var _ Node = (DeclSections)(nil)

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

var _ DeclSection = (*LabelDeclSection)(nil)
var _ astcore.DeclNode = (*LabelDeclSection)(nil)

func (*LabelDeclSection) canBeDeclSection() {}
func (m *LabelDeclSection) Children() Nodes { return Nodes{m.LabelId} }
func (m *LabelDeclSection) ToDeclarations() astcore.Decls {
	return astcore.Decls{astcore.NewDeclaration(m.LabelId, m)}
}

type LabelId = Ident

func NewLabelId(ident *Ident) *LabelId { return ident }
