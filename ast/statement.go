package ast

// - CompoundStmt
//   ```
//   BEGIN StmtList END
//   ```
type CompoundStmt struct {
	Node
	*StmtList
}

func (m *CompoundStmt) Children() Nodes { return Nodes{m.StmtList} }

// - StmtList
//   ```
//   Statement ';'
//   ```
type StmtList struct {
	Node
	*Statement
}

func (m *StmtList) Children() Nodes { return Nodes{m.Statement} }

// - Statement
//   ```
//   [LabelId ':'] [SimpleStatement | StructStmt]
//   ```
type Statement struct {
	Node
	LabelId *LabelId
	Body    StatementBody
}

func (m *Statement) Children() Nodes { return Nodes{m.LabelId, m.Body} }

type StatementBody interface {
	Node
	isStatementBody()
}

type SimpleStatement interface {
	StatementBody
	isSimpleStatement()
}

type StructStmt interface {
	StatementBody
	isStructStmt()
}

// - SimpleStatement
//   (CallStatement)
//   ```
//   Designator ['(' [ExprList] ')']
//   ```
//   (AssignStatement)
//   ```
//   Designator ':=' Expression
//   ```
//   (InheritedStatement)
//   ```
//   INHERITED
//   ```
//   (GotoStatement)
//   ```
//   GOTO LabelId
//   ```
// - StructStmt
//   ```
//   CompoundStmt
//   ```
//   ```
//   ConditionalStmt
//   ```
//   ```
//   LoopStmt
//   ```
//   ```
//   WithStmt
//   ```
//   ```
//   TryExceptStmt
//   ```
//   ```
//   TryFinallyStmt
//   ```
//   ```
//   RaiseStmt
//   ```
//   ```
//   AssemblerStmt
//   ```

type DesignatorStatement interface {
	SimpleStatement
	isDesignatorStatement()
}

//   (CallStatement)
//   ```
//   Designator ['(' [ExprList] ')']
//   ```
type CallStatement struct {
	DesignatorStatement
	Designator *Designator
	ExprList   ExprList // nil able
}

func (*CallStatement) isStatementBody()         {}
func (*CallStatement) isSimpleStatement()       {}
func (m *CallStatement) isDesignatorStatement() {}
func (m *CallStatement) Children() Nodes        { return Nodes{m.Designator, m.ExprList} }

//   (AssignStatement)
//   ```
//   Designator ':=' Expression
//   ```
type AssignStatement struct {
	DesignatorStatement
	Designator *Designator
	Expression *Expression
}

func (*AssignStatement) isStatementBody()         {}
func (*AssignStatement) isSimpleStatement()       {}
func (m *AssignStatement) isDesignatorStatement() {}
func (m *AssignStatement) Children() Nodes        { return Nodes{m.Designator, m.Expression} }
