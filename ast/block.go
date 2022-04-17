package ast

// - Block
//   ```
//   [DeclSection]
//   [ExportsStmt]...
//   CompoundStmt
//   [ExportsStmt]...
//   ```
type Block struct {
	DeclSection   *DeclSection
	ExportsStmts1 ExportsStmts
	CompoundStmt  *CompoundStmt
	ExportsStmts2 ExportsStmts
}

// - ExportsStmt
//   ```
//   EXPORTS ExportsItem [, ExportsItem]...
//   ```
type ExportsStmt struct {
	ExportsItems []*ExportsItem
}

type ExportsStmts []*ExportsStmt

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
	canBeDeclSection()
}

// - LabelDeclSection
//   ```
//   LABEL LabelId
//   ```
type LabelDeclSection struct {
	*LabelId
}

func (*LabelDeclSection) canBeDeclSection() {}

type LabelId = Ident
