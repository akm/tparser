package ast

// - CompoundStmt
//   ```
//   BEGIN StmtList END
//   ```
type CompoundStmt struct {
	*StmtList
}

// - StmtList
//   ```
//   Statement ';'
//   ```
type StmtList struct {
	*Statement
}

// - Statement
//   ```
//   [LabelId ':'] [SimpleStatement | StructStmt]
//   ```
type Statement struct {
	LabelId *LabelId
	Body    StatementBody
}

type StatementBody interface {
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
	GetDesignator() *Designator
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

func (*CallStatement) isStatementBody()             {}
func (*CallStatement) isSimpleStatement()           {}
func (m *CallStatement) GetDesignator() *Designator { return m.Designator }

//   (AssignStatement)
//   ```
//   Designator ':=' Expression
//   ```
type AssignStatement struct {
	DesignatorStatement
	Designator *Designator
	Expression *Expression
}

func (*AssignStatement) isStatementBody()             {}
func (*AssignStatement) isSimpleStatement()           {}
func (m *AssignStatement) GetDesignator() *Designator { return m.Designator }
