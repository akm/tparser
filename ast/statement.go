package ast

import "github.com/akm/tparser/ast/astcore"

// - CompoundStmt
//   ```
//   BEGIN StmtList END
//   ```
type CompoundStmt struct {
	StructStmt
	StmtList
}

func (*CompoundStmt) isStructStmt()     {}
func (m *CompoundStmt) Children() Nodes { return Nodes{m.StmtList} }

// - StmtList
//   ```
//   (Statement ';') ...
//   ```
type StmtList []*Statement // must implement Node
func (s StmtList) Children() Nodes {
	r := make(Nodes, len(s))
	for idx, i := range s {
		r[idx] = i
	}
	return r
}

// - Statement
//   ```
//   [LabelId ':'] [SimpleStatement | StructStmt]
//   ```
type Statement struct {
	Node
	LabelId *LabelId
	Body    StatementBody
}

func (m *Statement) Children() Nodes {
	res := Nodes{}
	if m.LabelId != nil {
		res = append(res, m.LabelId)
	}
	res = append(res, m.Body)
	return res
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

// DesignatorStatement is an interface for CallStatement and AssignStatement
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
func (m *CallStatement) Children() Nodes {
	res := Nodes{m.Designator}
	if m.ExprList != nil {
		res = append(res, m.ExprList)
	}
	return res
}

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
func (m *AssignStatement) Children() Nodes {
	return Nodes{m.Designator, m.Expression}
}

//   (InheritedStatement)
//   ```
//   INHERITED
//   ```
type InheritedStatement struct {
	SimpleStatement
	Ref *astcore.Declaration // reference to the ancestor method
}

func (*InheritedStatement) isStatementBody()   {}
func (*InheritedStatement) isSimpleStatement() {}
func (*InheritedStatement) Children() Nodes    { return Nodes{} }

//   (GotoStatement)
//   ```
//   GOTO LabelId
//   ```
type GotoStatement struct {
	SimpleStatement
	LabelId *LabelId
	Ref     *astcore.Declaration
}

func (*GotoStatement) isStatementBody()   {}
func (*GotoStatement) isSimpleStatement() {}
func (m *GotoStatement) Children() Nodes  { return Nodes{m.LabelId} }
